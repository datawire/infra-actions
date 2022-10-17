'use strict';

const core = require('@actions/core')
const github = require('@actions/github')
const fs = require('fs')

const kubeception = require('./lib/kubeception.js')
const { v4: uuidv4 } = require('uuid')
const utils = require('./lib/utils.js')

const registry = require('./registry.js')

async function create() {
  // inputs are defined in action metadata file
  const distribution = core.getInput('distribution')
  const action = core.getInput('action')
  const version = core.getInput('version')
  const kubeconfigPath = core.getInput('kubeconfig')

  if (action === 'expire') {
    return
  }

  switch (distribution.toLowerCase()) {
  case "kubeception":
    const clusterName = utils.getUniqueClusterName()
    core.exportVariable('clusterName', clusterName)

    core.notice(`Creating ${distribution} ${version} and writing kubeconfig to file: ${kubeconfigPath}!`)
    const kubeConfig = kubeception.createKluster(clusterName, version)
    kubeConfig.then(contents => { utils.writeKubeconfig(kubeconfig, contents) })
    break
  default:
    let provider = registry.getProvider(distribution)

    core.notice(`Creating ${distribution} ${version} and writing kubeconfig to file: ${kubeconfigPath}!`)
    let cluster = await provider.allocateCluster()
    core.saveState(registry.CLUSTER_NAME, cluster.name)

    core.notice(`Created ${distribution} cluster ${cluster.name}!`)

    let kubeconfig = await provider.makeKubeconfig(cluster)
    let contents = JSON.stringify(kubeconfig, undefined, 2) + "\n"
    utils.writeKubeconfig(kubeconfigPath, contents)

    core.notice(`Exporting KUBECONFIG as ${kubeconfigPath}`)
    core.exportVariable("KUBECONFIG", kubeconfigPath)
    break
  }
}

create().catch((error)=>{
  core.setFailed(error.message)
})
