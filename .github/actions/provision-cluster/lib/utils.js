"use strict";

const fs = require("fs");
const crypto = require("crypto");
const core = require("@actions/core");

function getUniqueClusterName(maxNameLength) {
  const repoName = process.env["GITHUB_REPOSITORY"].replace(/^.*\//, "");
  const branch = process.env["GITHUB_HEAD_REF"];
  const sha = process.env["GITHUB_SHA"].substring(0, 8);

  let name = `ci-${uid()}-${repoName}-${sha}-${branch}`;
  let sanitizedName = name
    .replace(/[^A-Za-z0-9-]/g, "-")
    .replace(/-+$/g, "")
    .toLowerCase()
    .substring(0, maxNameLength);

  if (sanitizedName.endsWith("-")) {
    sanitizedName = sanitizedName.substring(0, sanitizedName.length - 1);
  }

  return sanitizedName;
}

function writeFile(path, contents) {
  fs.writeFile(path, contents, (err) => {
    if (err) {
      core.setFailed(`${err}`);
    }
  });
}

// Convenience for sleeping in async functions/methods.
function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

// Construct a thunk that returns a fibonacci sequence using the supplied initial and max delays.
function fibonacciDelaySequence(initialDelay, maxDelay) {
  let prevFibonacciDelay = 0;
  let curFibonacciDelay = initialDelay;

  return () => {
    const result = curFibonacciDelay + prevFibonacciDelay;
    prevFibonacciDelay = curFibonacciDelay;
    curFibonacciDelay = result;
    if (typeof maxDelay === typeof undefined) {
      return result;
    } else {
      return Math.min(result, maxDelay);
    }
  };
}

// Construct a unique id.
function uid() {
  return crypto.randomBytes(16).toString("hex");
}

class Retry extends Error {}

class Transient extends Error {}

// Retry the supplied action with a fibonacci backoff until it returns true or timeout seconds have
// passed. Use minDelay and maxDelay to tune the delay times. The action should signal retry with
// `throw new Transient(...)` or `throw new Retry(...)`, and return upon success.
// The result of the final successful invocation will be returned.
async function fibonacciRetry(
  action,
  timeout = 600000,
  minDelay = 1000,
  maxDelay = 30000
) {
  let start = Date.now();
  let nextDelay = fibonacciDelaySequence(minDelay, maxDelay);

  let count = 0;
  let timeoutReached = false;

  do {
    count++;

    try {
      return await action();
    } catch (e) {
      if (!(e instanceof Transient) && !(e instanceof Retry)) {
        throw e;
      }
      let delay = nextDelay();
      let elapsed = Date.now() - start;
      let remaining = timeout - elapsed;
      if (remaining > 0) {
        let t = Math.min(delay, remaining);
        core.info(`Error (${e.message}) retrying after ${t / 1000}s ...`);
        await sleep(t);
      } else {
        timeoutReached = true;
      }

      if (timeoutReached) {
        throw new Error(
          `Error (${e.message}) failing after ${count} attempts over ${
            elapsed / 1000
          }s.`
        );
      }
    }
  } while (!timeoutReached);
}

module.exports = {
  getUniqueClusterName,
  writeFile,
  sleep,
  fibonacciDelaySequence,
  uid,
  Retry,
  Transient,
  fibonacciRetry,
};
