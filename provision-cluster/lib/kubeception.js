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
  if (!name) {
    throw new Error('Function createKluster() needs a Kluster name');
  }

  if (!version) {
    throw Error('Function createKluster() needs a Kluster version');
  }

  const kubeceptionToken = core.getInput('kubeceptionToken');
  if (!kubeceptionToken) {
    throw Error(`kubeceptionToken is missing. Make sure that input parameter kubeceptionToken was provided`);
  }

  const client = getHttpClient();

  const oneDay = 86400
  const url = `https://sw.bakerstreet.io/kubeception/api/klusters/${name}?version=${version}&wait=true&timeoutSecs=${oneDay}&EnableSNIRelay=true`
  let response = await client.put(url);
  if (!response || !response.message) {
    throw Error("Unknown error getting response");
  }

  //Temporarily retry 504 errors to be able to use this in a sample job
  if (response.message.statusCode == 504) {
    core.warning("Retrying 504 error")
    response = await client.put(url)
  }

  //Temporarily retry 504 errors to be able to use this in a sample job
  if (response.message.statusCode == 504) {
    core.warning("Retrying 504 error")
    response = await client.put(url)
  }

  if (response.message.statusCode != 200) {
    throw Error(`Expected status code 200 but got ${response.message.statusCode}`);
  }

  const body = await response.readBody();

  return body;
};

async function deleteKluster(name) {
  if (!name) {
    throw Error('Function deleteKluster() needs a Kluster name');
  }

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