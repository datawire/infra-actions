const { v4: uuidv4 } = require('uuid');
const fs = require('fs');

function getUniqueClusterName(maxNameLength) {
  const repoName = process.env['GITHUB_REPOSITORY'].replace(/^.*\//, '');
  const branch = process.env['GITHUB_HEAD_REF'];
  const sha = process.env['GITHUB_SHA'].substring(0, 8);
	const uuid = uuidv4().replace('-', '').substring(0, 8);

  let name = `test-${uuid}-${repoName}-${sha}-${branch}`;
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

module.exports = { getUniqueClusterName, writeFile};