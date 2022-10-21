'use strict';

const core = require('@actions/core')
const github = require('@actions/github')

const registry = require('./registry.js')

async function do_delete() {
  // inputs are defined in action metadata file
  const distribution = core.getInput('distribution')
  let clusterName = core.getState(registry.CLUSTER_NAME)

  let provider = registry.getProvider(distribution)

  let promises = []
  promises.push(expire(provider))

  if (typeof clusterName !== typeof undefined) {
    core.notice(`Deleting ${distribution} cluster ${clusterName}!`)
    promises.push(delete_allocated(provider, clusterName))
  }

  Promise.allSettled(promises)
}

async function expire(provider) {
  return provider.expireClusters()
}

async function delete_allocated(provider, name) {
  let cluster = await provider.getCluster(clusterName)
  return provider.deleteCluster(cluster)
}

do_delete().catch((error)=>{
  core.setFailed(error.message)
})
