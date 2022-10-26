'use strict';

const container = require('@google-cloud/container')
const crypto = require('crypto')
const utils = require('./lib/utils.js')

const STATUS_ENUM = container.protos.google.container.v1.Operation.Status

// Every cluster created by this action is labeled with provisioned-category=ephemeral
const CATEGORY_PROPERTY = 'provisioned-category'
// Every cluster created by this action will have a provisioned-lifespan that defines when it is
// ok to delete the cluster.
const LIFESPAN_PROPERTY = 'provisioned-lifespan'
const EPHEMERAL = 'ephemeral'
// Default lifespan of 30 minutes.
const DEFAULT_LIFESPAN = 1800 // 30 minutes
//const DEFAULT_LIFESPAN = 600 // 10 minutes for dev

// The Client class is a convenience wrapper around the google API that allows for sharing of some
// of the boilerplate between different operations.
class Client {

  constructor(zone, gkeClient) {
    if (typeof gkeClient == typeof undefined) {
      gkeClient = new container.v1.ClusterManagerClient()
    }
    this.client = gkeClient
    this.project = null
    this.zone = zone
  }

  async getProjectId() {
    if (this.project === null) {
      this.project = await this.client.getProjectId()
    }
    return this.project
  }

  // Compute the location used in multiple methods.
  async getLocation() {
    return `projects/${await this.getProjectId()}/locations/${this.zone}`
  }

  // Create a new cluster with a unique name, wait for it to be fully provisioned, and then fetch
  // and return the resulting cluster object.
  async allocateCluster() {
    let name = `test-${utils.uid()}`
    let cluster = {
      name: name,
      network: 'default',
      initialNodeCount: 1,
      nodeConfig: {
        machineType: 'e2-standard-2',
      }
    }
    await this.createCluster(cluster)
    return this.getCluster(name)
  }

  // Get a cluster by name.
  async getCluster(name) {
    const [cluster] = await this.client.getCluster({name: `${await this.getLocation()}/clusters/${name}`})
    return cluster
  }

  // Make a functioning kubeconfig from a cluster object.
  async makeKubeconfig(cluster) {
    let token = await this.client.auth.getAccessToken()

    let kubeconfig = {
      apiVersion: "v1",
      kind: "Config",
      clusters: [{
        cluster: {
          "certificate-authority-data": cluster.masterAuth.clusterCaCertificate,
          server: `https://${cluster.endpoint}`
        },
        name: "gke-cluster"
      }],
      users: [{name: "gke-user", user: {token: token}}],
      contexts: [{
        context: {
          cluster: "gke-cluster",
          namespace: "default",
          user: "gke-user"
        },
        name: "gke-context"
      }],
      "current-context": "gke-context"
    }

    return kubeconfig
  }

  // Iterate over all the clusters in the zone and delete any expired clusters.
  async expireClusters(lifespanOverride) {
    let promises = []
    for (let c of await this.listClusters()) {
      promises.push(this.maybeExpireCluster(c, lifespanOverride))
    }
    return Promise.allSettled(promises)
  }

  async listClusters() {
    const [response] = await this.client.listClusters({
      projectId: await this.getProjectId(),
      zone: this.zone,
    })
    return response.clusters
  }

  // Create the supplied cluster. This method will automatically add labels to mark the cluster as
  // having been created by this action and it will provide a default lifespan label of 30 minutes.
  async createCluster(cluster) {
    if (typeof cluster.resourceLabels === typeof undefined) {
      cluster.resourceLabels = {}
    }
    if (typeof cluster.resourceLabels[CATEGORY_PROPERTY] === typeof undefined) {
      cluster.resourceLabels[CATEGORY_PROPERTY] = EPHEMERAL
    }
    if (typeof cluster.resourceLabels[LIFESPAN_PROPERTY] === typeof undefined) {
      cluster.resourceLabels[LIFESPAN_PROPERTY] = DEFAULT_LIFESPAN
    }

    const [operation] = await this.client.createCluster({parent: await this.getLocation(), cluster: cluster})
    return this.awaitOperation(operation)
  }

  // Delete the given cluster. This method will throw an exception if the supplied cluster does not
  // have the appropriate labels that indicate the cluster was created by this github action. Pass
  // in force=true to override this check.
  async deleteCluster(cluster, force=false) {
    let name = cluster.name

    if (!force && cluster.resourceLabels[CATEGORY_PROPERTY] !== EPHEMERAL) {
      return new Operation(false, `Cannot delete cluster ${name}, it is not ephemeral.`)
    }

    try {
      const [op] = await this.client.deleteCluster({name: `${await this.getLocation()}/clusters/${name}`})
      return Operation.wrap(op)
    } catch (error) {
      return new Operation(false, `Error deleting cluster: ${error}`)
    }
  }

  // Check if the cluster is both ephemeral and old enough that it should be deleted, and if it is,
  // then delete it. Please note that it is *really* important that this code does not delete GKE
  // clusters that were provisioned by hand or any means other than this action. That is why we
  // check for both the `provisioned-category` and `provisioned-lifespan` labels and ignore the
  // cluster if they are not present and set to the correct value. We are expecting that only
  // clusters provisioned by this action will have those labels.
  async maybeExpireCluster(cluster, lifespanOverride) {
    let labels = cluster.resourceLabels

    let category = labels[CATEGORY_PROPERTY]
    if (category !== EPHEMERAL) {
      console.log(`Ignoring cluster ${cluster.name} because it has not ephemeral.`)
      return
    }

    // Lifespan is in seconds
    let lifespan = labels[LIFESPAN_PROPERTY]
    if (typeof lifespan === typeof undefined) {
      console.log(`Keeping cluster ${cluster.name} because it has no provisioned-lifespan label.`)
      return
    }

    let lifespanMillis = 0

    if (typeof lifespanOverride !== typeof undefined) {
      lifespanMillis = 1000*lifespanOverride
    } else {
      if (typeof lifespan === "string" && lifespan !== "") {
        lifespanMillis = Number(lifespan)*1000
      }
    }

    if (lifespanMillis <= 0) {
      console.log(`Keeping ${cluster.name} because the lifespan is <= 0.`)
      return
    }

    let ageMillis = Date.now() - Date.parse(cluster.createTime)
    if (ageMillis < lifespanMillis) {
      console.log(`Keeping ${cluster.name} because ${ageMillis/1000}s < ${lifespanMillis/1000}s.`)
      return
    }

    console.log(`Deleting ${cluster.name} because ${ageMillis/1000}s >= ${lifespanMillis/1000}s.`)
    let op = await this.deleteCluster(cluster)
    console.log(op.status)
  }

  // Wait for the supplied operation to finish by polling up to limit times.
  async awaitOperation(operation) {
    utils.fibonacciRetry(async ()=>{
      const op = await this.getOperation(operation)
      if (op.done) {
        return op
      } else {
        throw new utils.Transient(op.status)
      }
    })
  }

  // Get the current status of the supplied operation.
  async getOperation(operation) {
    let location = await this.getLocation()
    let opId = `${location}/operations/${operation.name}`
    const [op] = await this.client.getOperation({name: opId})
    return Operation.wrap(op)
  }

}

// The Operation object is used to report the status of a long running procedures.
class Operation {

  // Wrap the google supplied operation class and convert it into something simpler.
  static wrap(op) {
    const done = op.status == STATUS_ENUM[STATUS_ENUM.DONE]

    let url = op.targetLink
    let idx = url.lastIndexOf('/')
    let name = url
    if (idx >= 0) {
      name = url.substring(idx+1)
    }
    if (op.detail) {
      return new Operation(done, `${op.operationType} ${name} ${op.status}: ${op.detail}`)
    } else {
      return new Operation(done, `${op.operationType} ${name} ${op.status}`)
    }
  }

  constructor(done, status) {
    this.done = done
    this.status = status
  }

}

module.exports = {
  Client
}
