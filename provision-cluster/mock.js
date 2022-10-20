const gke = require('./gke.js')

class MockGKE {
  constructor() {
    this.auth = new MockAuth()
  }
}

const MOCK_ACCESS_TOKEN = "mock-access-token"
const MOCK_ZONE = 'mock-zone'

class MockAuth {

  getAccessToken() {
    return MOCK_ACCESS_TOKEN
  }

}

// Return a client but with a mocked GKE client underneath for testing purposes.
function Client() {
  let client = new gke.Client(MOCK_ZONE)
  client.client = new MockGKE()
  return client
}

module.exports = {
  Client,
  MOCK_ACCESS_TOKEN,
  MOCK_ZONE
}
