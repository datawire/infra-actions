'use strict'
const core = require('@actions/core')
const httpClient = require('@actions/http-client')
const httpClientLib = require('@actions/http-client/lib/auth.js')
const errorHandling = require('./errors.js')

function getHttpClient() {
  const userAgent = 'datawire/provision-cluster'

  const kubeceptionToken = core.getInput('kubeceptionToken')
  if (!kubeceptionToken) {
    throw Error(`kubeceptionToken is missing. Make sure that input parameter kubeceptionToken was provided`)
  }

  const credentialHandler = new httpClientLib.BearerCredentialHandler(kubeceptionToken)
  return new httpClient.HttpClient(userAgent, [credentialHandler])
}

async function createKluster(name, version) {
  if (!name) {
    throw new Error('Function createKluster() needs a Kluster name')
  }

  if (!version) {
    throw Error('Function createKluster() needs a Kluster version')
  }

  const kubeceptionToken = core.getInput('kubeceptionToken')
  if (!kubeceptionToken) {
    throw Error(`kubeceptionToken is missing. Make sure that input parameter kubeceptionToken was provided`)
  }

  const client = getHttpClient()

  const oneDay = 86400
  const response = await client.put(`https://sw.bakerstreet.io/kubeception/api/klusters/${name}?version=${version}&timeoutSecs=${oneDay}`)
  if (!response || !response.message) {
    throw Error("Unknown error getting response")
  }

  if (response.message.statusCode == 425) {
    throw new errorHandling.TransientError('Kluster is not ready')
  }

  if (response.message.statusCode != 200) {
    throw Error(`Expected status code 200 but got ${response.message.statusCode}`)
  }

  const body = await response.readBody()

  return body
}

async function deleteKluster(name) {
  if (!name) {
    throw Error('Function deleteKluster() needs a Kluster name')
  }

  const client = getHttpClient()

  const response = await client.del(`https://sw.bakerstreet.io/kubeception/api/klusters/${name}`)
  if (!response || !response.message) {
    throw Error("Unknown error getting response")
  }

  if (response.message.statusCode != 200) {
    throw Error(`Expected status code 200 but got ${response.message.statusCode}`)
  }
}

module.exports = { createKluster, deleteKluster}