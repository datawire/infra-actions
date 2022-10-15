const crypto = require('crypto')

function getUniqueClusterName() {
  const repoName = process.env['GITHUB_REPOSITORY'].replace(/^.*\//, '');
  const branch = process.env['GITHUB_HEAD_REF'];
  const sha = process.env['GITHUB_SHA'].substring(0, 8);

  let name = `ci-${uniqueId()}-${repoName}-${sha}-${branch}`;
  let sanitizedName = name.replace(/[^A-Za-z0-9-]/g, '-').replace(/-+$/g, '').toLowerCase().substring(0, 63);

	return sanitizedName;
}

function uniqueId() {
  return crypto.randomBytes(16).toString("hex")
}

module.exports = { getUniqueClusterName};