const core = require('@actions/core')
const utils = require('./utils.js')
const kubeception = require('./kubeception.js');

function runWithRetry(promise) {
  core.info(`executing ${promise}`)
  result = promise
    .then(result => {
      core.debug(`Finished creation`)
      return result
      })
    .catch(error => {
      core.debug(`Caught temporary error ${error}`)
      setTimeout(promise, 1000)
     })

  return result
}

module.exports = { runWithRetry }