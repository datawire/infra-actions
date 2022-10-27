'use strict';

const gke = require('./gke.js')
const fs = require('fs')
const path = require('path')

const container = require('@google-cloud/container')
const STATUS_ENUM = container.protos.google.container.v1.Operation.Status

const MOCK = {
  PROJECT_ID: 'mock-project-id',
  ACCESS_TOKEN: 'mock-access-token',
  ZONE: 'mock-zone',
  OPERATION_NAME: 'mock-operation-name',
  OPERATION_TARGET_LINK: 'https://mock/mock-cluster'
}

const cluster = JSON.parse(fs.readFileSync(path.join(__dirname, 'cluster.json')))

class MockGKE {
  constructor() {
    this.auth = new MockAuth()
  }

  async getProjectId() {
    return MOCK.PROJECT_ID
  }

  async createCluster() {
    return [new MockOp()]
  }

  async deleteCluster() {
    return [new MockOp()]
  }

  async getCluster() {
    return [cluster]
  }

  async getOperation(op) {
    return [new MockOp()]
  }

}

class MockAuth {
  
  getAccessToken() {
    return MOCK.ACCESS_TOKEN
  }

}

class MockOp {
  constructor() {
    this.name = MOCK.OPERATION_NAME
    this.targetLink = cluster.selfLink
    this.status = STATUS_ENUM[STATUS_ENUM.DONE]
  }
}

// Return a client but with a mocked GKE client underneath for testing purposes.
function Client() {
  let client = new gke.Client(MOCK.ZONE, new MockGKE())
  return client
}

module.exports = {
  MOCK,
  Client,
  MockGKE
}
