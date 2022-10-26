const fs = require('fs');
const crypto = require('crypto')
const core = require('@actions/core')

function getUniqueClusterName(maxNameLength) {
  const repoName = process.env['GITHUB_REPOSITORY'].replace(/^.*\//, '');
  const branch = process.env['GITHUB_HEAD_REF'];
  const sha = process.env['GITHUB_SHA'].substring(0, 8);

  let name = `ci-${uniqueId()}-${repoName}-${sha}-${branch}`;
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

function uniqueId() {
  return crypto.randomBytes(16).toString("hex")
}

module.exports = { getUniqueClusterName, writeFile};
