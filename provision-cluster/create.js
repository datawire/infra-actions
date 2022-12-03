'use strict';

const core = require('@actions/core')
const github = require('@actions/github')

const registry = require('./lib/registry.js')
const utils = require('./lib/utils.js')

async function create() {
  // inputs are defined in action metadata file
  const distribution = core.getInput('distribution')
  const action = core.getInput('action')
  const version = core.getInput('version')
  const lifespan = core.getInput('lifespan')
  const kubeconfigPath = core.getInput('kubeconfig')

  let provider = registry.getProvider(distribution)

  if (action === 'expire') {
    return
  }

  core.notice(`Creating ${distribution} ${version} and writing kubeconfig to file: ${kubeconfigPath}!`)
  let cluster = await provider.allocateCluster(version, lifespan)
  core.saveState(registry.CLUSTER_NAME, cluster.name)

  core.setOutput("clusterName", cluster?.name);
  core.setOutput("projectId", cluster?.project);
  core.setOutput("location", cluster?.zone);

  core.notice(`Creating ${distribution} cluster ${cluster.name} ...`)
  let kubeconfig = await provider.makeKubeconfig(cluster)
  core.notice(`Cluster created: ${cluster.name}!`)
  let contents = JSON.stringify(kubeconfig, undefined, 2) + "\n"
  utils.writeFile(kubeconfigPath, contents)

  core.notice(`Kubeconfig written to ${kubeconfigPath}.`)
  core.notice(`Exporting KUBECONFIG as ${kubeconfigPath}`)
  core.exportVariable("KUBECONFIG", kubeconfigPath)
}

create().catch((error)=>{
  core.setFailed(error.message)
})
