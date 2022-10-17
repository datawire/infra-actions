'use strict';
const core = require('@actions/core');
const github = require('@actions/github');
const kubeception = require('./lib/kubeception.js');
const { v4: uuidv4 } = require('uuid');
const utils = require('./lib/utils.js');

try {
  // inputs are defined in action metadata file
  const distribution = core.getInput('distribution');
  const version = core.getInput('version');
  const kubeconfig = core.getInput('kubeconfig');

  const clusterName = utils.getUniqueClusterName()
  core.exportVariable('clusterName', clusterName);

  switch(distribution.toLowerCase()) {
    case "kubeception": {
      const kubeConfig = kubeception.createKluster(clusterName, version);
      kubeConfig.then(contents => { utils.writeFile(kubeconfig, contents); });
      break;
    }
    default: {
      console.log(`Creating ${distribution} ${version} and writing kubeconfig to file: ${kubeconfig}!`);
      let kubeconfigContents = `Mock kubeconfig file for ${distribution} ${version}.\n`;
      utils.writeFile(kubeconfig, kubeconfigContents);
      break;
    }
  }

  // Get the JSON webhook payload for the event that triggered the workflow
  //const payload = JSON.stringify(github.context.payload, undefined, 2)
  //console.log(`The event payload: ${payload}`);
} catch (error) {
  core.setFailed(error.message);
}

