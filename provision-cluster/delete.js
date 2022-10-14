'use strict';
const core = require('@actions/core');
const github = require('@actions/github');
const kubeception = require('./kubeception.js');

try {
  // inputs are defined in action metadata file
  const distribution = core.getInput('distribution');
  const version = core.getInput('version');
  // Get the JSON webhook payload for the event that triggered the workflow
  //const payload = JSON.stringify(github.context.payload, undefined, 2)
  //console.log(`The event payload: ${payload}`);

  switch(distribution.toLowerCase()) {
    case "kubeception": {
      kubeception.deleteKluster("aosorio-test-kluster").then(
        console.log(`Deleting ${distribution} ${version}!`)
      );
      break;
    }
    default: {
      console.log(`Deleting ${distribution} ${version}!`);
      break;
    }
  }
} catch (error) {
  core.setFailed(error.message);
}
