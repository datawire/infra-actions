'use strict';
const core = require('@actions/core');
const github = require('@actions/github');
const fs = require('fs');
const kubeception = require('./kubeception.js');

try {
  // inputs are defined in action metadata file
  const distribution = core.getInput('distribution');
  const version = core.getInput('version');
  const kubeconfig = core.getInput('kubeconfig');

  switch(distribution.toLowerCase()) {
    case "kubeception": {
      const kubeConfig = kubeception.createKluster("aosorio-test-kluster", version);
      kubeConfig.then(contents => { writeKubeconfig(kubeconfig, contents); });
      break;
    }
    default: {
      console.log(`Creating ${distribution} ${version} and writing kubeconfig to file: ${kubeconfig}!`);
      let kubeconfigContents = `Mock kubeconfig file for ${distribution} ${version}.\n`;
      writeKubeconfig(kubeconfig, kubeconfigContents);
      break;
    }
  }

  // Get the JSON webhook payload for the event that triggered the workflow
  //const payload = JSON.stringify(github.context.payload, undefined, 2)
  //console.log(`The event payload: ${payload}`);
} catch (error) {
  core.setFailed(error.message);
}

function writeKubeconfig(path, contents) {
  fs.writeFile(path, contents, err => {
    if (err) {
      core.setFailed(`${err}`);
    }
  });
}