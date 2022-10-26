'use strict';
const core = require('@actions/core');
const httpClient = require('@actions/http-client');
const httpClientLib = require('@actions/http-client/lib/auth.js');
const utils = require('./utils.js')
const yaml = require('yaml')

const MAX_KLUSTER_NAME_LEN = 63
const oneHour = 60*1000

class Client {

  async allocateCluster(version) {
    const clusterName = utils.getUniqueClusterName(MAX_KLUSTER_NAME_LEN)
    const kubeConfig = await createKluster(clusterName, version)
    return {
      "name": clusterName,
      "config": kubeConfig
    }
  }

  async makeKubeconfig(cluster) {
    return yaml.parse(cluster.config)
  }

  async getCluster(clusterName) {
    return clusterName
  }

  async deleteCluster(clusterName) {
    return deleteKluster(clusterName)
  }

  async expireClusters() {
    // Kubeception automatically expires klusters, no client side expiration is required.
  }
}

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

  return utils.fibonacciRetry(async ()=>{
    const response = await client.put(`https://sw.bakerstreet.io/kubeception/api/klusters/${name}?version=${version}&timeoutSecs=${oneHour}`)
    if (!response || !response.message) {
      throw Error("Unknown error getting response")
    }

    if (response.message.statusCode == 200) {
      return await response.readBody()
    } else if (response.message.statusCode == 425) {
      throw new utils.Transient(`status code ${response.message.statusCode}`)
    } else {
      let body = await response.readBody()
      throw new Error(`Status code ${response.message.statusCode}: ${body}`)
    }
  })
}

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

module.exports = { Client, createKluster, deleteKluster};
