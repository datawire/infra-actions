'use strict';

const core = require('@actions/core')
const github = require('@actions/github')

const registry = require('./lib/registry.js')
const slack = require('./lib/slack.js')

async function do_delete() {
  // inputs are defined in action metadata file
  const distribution = core.getInput('distribution')
  const clusterName = core.getState(registry.CLUSTER_NAME)

  let provider = registry.getProvider(distribution)

  let promises = []
  promises.push(expire(provider, distribution))

  if (typeof clusterName !== typeof undefined && clusterName !== "") {
    core.notice(`Deleting ${distribution} cluster ${clusterName}!`)
    promises.push(delete_allocated(provider, clusterName))
  }

  return Promise.all(promises)
}

async function expire(provider, distribution) {
  let orphaned = await provider.expireClusters()

  if (orphaned.length == 0) {
    return
  }

  core.notice(`Orhpaned Clusters: ${orphaned.join(', ')}`)
  slack.notify(`Orphaned clusters:\n\n - ${orphaned.join("\n - ")}`)
}

async function delete_allocated(provider, name) {
  let cluster = await provider.getCluster(name)
  return provider.deleteCluster(cluster)
}

do_delete().catch((error)=>{
  core.setFailed(error.message)
})
