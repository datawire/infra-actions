'use strict';
const core = require('@actions/core');
const github = require('@actions/github');
const fs = require('fs');
const kubeception = require('./lib/kubeception.js');
const { v4: uuidv4 } = require('uuid');

try {
  // inputs are defined in action metadata file
  const distribution = core.getInput('distribution');
  const version = core.getInput('version');
  const kubeconfig = core.getInput('kubeconfig');

	let clusterName = getUniqueClusterName()
  core.exportVariable('clusterName', clusterName);

	switch(distribution) {
	   case "Kubeception": {
	      const kubeConfig = kubeception.createKluster(clusterName, version);
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

function getUniqueClusterName() {
  const repoName = process.env['GITHUB_REPOSITORY'].replace(/^.*\//, '');
  const branch = process.env['GITHUB_HEAD_REF'];
  const sha = process.env['GITHUB_SHA'].substring(0, 8);
	const uuid = uuidv4().replace('-', '').substring(0, 8);

  let name = `${uuid}-${repoName}-${sha}-${branch}`;
  let sanitizedName = name.replace(/[^A-Za-z0-9-]/g, '-').toLowerCase().substring(0, 63);

	return sanitizedName;
}