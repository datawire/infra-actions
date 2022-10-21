'use strict';

const gke = require('./gke.js')

const CLUSTER_NAME = 'CLUSTER_NAME'

const distributions = {
  "gke": new gke.Client()
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
