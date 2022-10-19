'use strict';
const core = require('@actions/core');
const github = require('@actions/github');
const kubeception = require('./lib/kubeception.js');
const utils = require('./lib/utils.js');
const actionErrors = require('./lib/error-handling.js')

const MAX_KLUSTER_NAME_LEN = 63

try {
  // inputs are defined in action metadata file
  const distribution = core.getInput('distribution');
  const version = core.getInput('version');
  const kubeconfig = core.getInput('kubeconfig');

//  const clusterName = utils.getUniqueClusterName(MAX_KLUSTER_NAME_LEN);
  const clusterName = 'test-aosorio';
  core.exportVariable('clusterName', clusterName);

  switch(distribution.toLowerCase()) {
    case "kubeception": {
      const promise = kubeception.createKluster(clusterName, version)
      const kubeConfig = actionErrors.runWithRetry(promise)
      kubeConfig.then(console.log)
      utils.writeFile(kubeconfig, "test");


//      kubeConfig.then(contents => { utils.writeFile(kubeconfig, contents); });
      core.debug(`Finished creating kluster`);

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
  core.setFailed(`Error creating cluster. ${error.message}`);
}

