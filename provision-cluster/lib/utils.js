const { v4: uuidv4 } = require('uuid');

function getUniqueClusterName() {
  const repoName = process.env['GITHUB_REPOSITORY'].replace(/^.*\//, '');
  const branch = process.env['GITHUB_HEAD_REF'];
  const sha = process.env['GITHUB_SHA'].substring(0, 8);
	const uuid = uuidv4().replace('-', '').substring(0, 8);

  let name = `${uuid}-${repoName}-${sha}-${branch}`;
  let sanitizedName = name.replace(/[^A-Za-z0-9-]/g, '-').toLowerCase().substring(0, 63);

	return sanitizedName;
}

module.exports = { getUniqueClusterName};