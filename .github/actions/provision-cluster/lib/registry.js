'use strict';

const gke = require('./gke.js')
const kubeception = require('./kubeception.js')

const CLUSTER_NAME = 'CLUSTER_NAME'
const clusterZone = 'us-central1-b'

const distributions = {
  "gke": new gke.Client(clusterZone),
  "kubeception": new kubeception.Client(kubeception.getHttpClient())
}

function getProvider(distribution) {
  let result = distributions[distribution.toLowerCase()]
  if (typeof result === typeof undefined) {
    throw new Error(`unknown distribution: ${distribution}`)
  }
  return result
}

module.exports = {
  getProvider,
  CLUSTER_NAME
}
