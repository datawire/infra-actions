'use strict';

const core = require('@actions/core')
const github = require('@actions/github')
const fs = require('fs')

const registry = require('./registry.js')

async function create() {
  // inputs are defined in action metadata file
  const distribution = core.getInput('distribution')
  const version = core.getInput('version')
  const kubeconfigPath = core.getInput('kubeconfig')

  let provider = registry.getProvider(distribution)

  core.notice(`Creating ${distribution} ${version} and writing kubeconfig to file: ${kubeconfigPath}!`)
  let cluster = await provider.allocateCluster()
  core.saveState(registry.CLUSTER_NAME, cluster.name)

  core.notice(`Created ${distribution} cluster ${cluster.name}!`)

  let kubeconfig = await provider.makeKubeconfig(cluster)
  let contents = JSON.stringify(kubeconfig, undefined, 2) + "\n"
  fs.writeFile(kubeconfigPath, contents, err => {
    if (err) {
      core.setFailed(`${err}`)
    }
  })

  core.notice(`Exporting KUBECONFIG as ${kubeconfigPath}`)
  core.exportVariable("KUBECONFIG", kubeconfigPath)
}

create().catch((error)=>{
  core.setFailed(error.message)
})
