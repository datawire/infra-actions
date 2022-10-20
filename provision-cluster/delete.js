'use strict';
const core = require('@actions/core');
const github = require('@actions/github');
const kubeception = require('./lib/kubeception.js');

try {
  // inputs are defined in action metadata file
  const distribution = core.getInput('distribution');
  const version = core.getInput('version');
  // Get the JSON webhook payload for the event that triggered the workflow
  //const payload = JSON.stringify(github.context.payload, undefined, 2)
  //console.log(`The event payload: ${payload}`);

  const clusterName = process.env['clusterName'];
  if (!clusterName) {
    throw Error(`Variable clusterName is undefined`);
  }

	switch(distribution) {
	   case "Kubeception": {
	      kubeception.deleteKluster(clusterName).then(
	        console.log(`Deleting Kluster ${clusterName}!`)
	      );
	      break;
	   }
	   default: {
			  console.log(`Deleting ${distribution} ${version}!`);
	      break;
	   }
	}
} catch (error) {
  console.log(`Error creating cluster ${error}`);
  core.setFailed(error.message);
}
