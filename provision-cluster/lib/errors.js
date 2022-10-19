const core = require('@actions/core')
const utils = require('./utils.js')
const kubeception = require('./kubeception.js');

// These errors will be caught by the retry logic and the action that failed will be executed again.
// Any action that throws this error should be idempotent.
class TransientError extends Error {
  constructor(message) {
    super(message)
    this.name = "TransientError"
  }
}

module.exports = { TransientError }