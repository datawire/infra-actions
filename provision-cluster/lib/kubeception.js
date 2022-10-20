'use strict';
const core = require('@actions/core');
const httpClient = require('@actions/http-client');
const httpClientLib = require('@actions/http-client/lib/auth.js');

function getHttpClient() {
	const userAgent = 'datawire/provision-cluster'

  const kubeceptionToken = core.getInput('kubeceptionToken');
  if (!kubeceptionToken) {
    throw Error(`kubeceptionToken is missing. Make sure that input parameter kubeceptionToken was provided`);
  }

	const credentialHandler = new httpClientLib.BearerCredentialHandler(kubeceptionToken);
  return new httpClient.HttpClient(userAgent, [credentialHandler]);
}

async function createKluster(name, version) {
  const kubeceptionToken = core.getInput('kubeceptionToken');
  if (!kubeceptionToken) {
    throw Error(`kubeceptionToken is missing. Make sure that input parameter kubeceptionToken was provided`);
  }

  const client = getHttpClient();

	const oneDay = 86400
	const response = await client.put(`https://sw.bakerstreet.io/kubeception/api/klusters/${name}?version=${version}&wait=true&timeoutSecs=${oneDay}`);
	if (!response || !response.message) {
		throw Error("Unknown error getting response");
	}

	if (response.message.statusCode != 200) {
		throw Error(`Expected status code 200 but got ${response.message.statusCode}`);
	}

	const body = await response.readBody();

	return body;
};

async function deleteKluster(name) {
  const client = getHttpClient();

	const response = await client.del(`https://sw.bakerstreet.io/kubeception/api/klusters/${name}`);
	if (!response || !response.message) {
		throw Error("Unknown error getting response");
	}

	if (response.message.statusCode != 200) {
		throw Error(`Expected status code 200 but got ${response.message.statusCode}`);
	}
}

module.exports = { createKluster, deleteKluster};