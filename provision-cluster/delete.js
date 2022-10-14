'use strict';

const core = require('@actions/core')
const github = require('@actions/github')

const kubeception = require('./kubeception.js')
const registry = require('./registry.js')

async function do_delete() {
  // inputs are defined in action metadata file
  const distribution = core.getInput('distribution')

  switch(distribution.toLowerCase()) {
  case "kubeception":
    kubeception.deleteKluster("aosorio-test-kluster").then(
      console.log(`Deleting ${distribution} ${version}!`)
    )
    break
  default:
    let clusterName = core.getState(registry.CLUSTER_NAME)

    let provider = registry.getProvider(distribution)

    let promises = []
    promises.push(expire(provider))

    if (typeof clusterName !== typeof undefined && clusterName !== "") {
      core.notice(`Deleting ${distribution} cluster ${clusterName}!`)
      promises.push(delete_allocated(provider, clusterName))
    }

    return Promise.all(promises)
  }
}

async function expire(provider) {
  return provider.expireClusters()
}

async function delete_allocated(provider, name) {
  let cluster = await provider.getCluster(name)
  return provider.deleteCluster(cluster)
}

do_delete().catch((error)=>{
  core.setFailed(error.message)
})
