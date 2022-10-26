const fs = require('fs');
const crypto = require('crypto')
const core = require('@actions/core')

function getUniqueClusterName(maxNameLength) {
  const repoName = process.env['GITHUB_REPOSITORY'].replace(/^.*\//, '');
  const branch = process.env['GITHUB_HEAD_REF'];
  const sha = process.env['GITHUB_SHA'].substring(0, 8);

  let name = `ci-${uid()}-${repoName}-${sha}-${branch}`;
  let sanitizedName = name.replace(/[^A-Za-z0-9-]/g, '-').replace(/-+$/g, '').toLowerCase().substring(0, maxNameLength);

	return sanitizedName;
}

function writeFile(path, contents) {
  fs.writeFile(path, contents, err => {
    if (err) {
      core.setFailed(`${err}`);
    }
  });
}

// Convenience for sleeping in async functions/methods.
function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

// Construct a thunk that returns a fibonacci sequence using the supplied initial and max delays.
function fibonacciDelaySequence(initialDelay, maxDelay) {
  let prevFibonacciDelay = 0
  let curFibonacciDelay = initialDelay

  return () => {
    const result = curFibonacciDelay + prevFibonacciDelay
    prevFibonacciDelay = curFibonacciDelay
    curFibonacciDelay = result
    if (typeof maxDelay === typeof undefined) {
      return result
    } else {
      return Math.min(result, maxDelay)
    }
  }
}

// Construct a unique id.
function uid() {
  return crypto.randomBytes(16).toString("hex")
}

module.exports = { getUniqueClusterName, writeFile, sleep, fibonacciDelaySequence, uid };
