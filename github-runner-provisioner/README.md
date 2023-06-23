# Runner Service

This service is based on the [echo
template](https://github.com/datawire/infrastructure/tree/master/echo). Please view the
[README](https://github.com/datawire/infrastructure/tree/master/echo) for details about the dev loop
and how it works.

# Architecture

We use the GitHub-Runner-Provisioner to serve a webhook to GitHub Actions. GitHub will send any 
Actions events to the GRP running in Skunkworks, which will parse those events looking for 
workflows that request special labels in their `runs-on` property. 

Using the GitHub Self-Hosted Runner binaries we then spin up the custom runners in one of our 
supported runner providers - currently AWS and CodeMagic. Supported runners are configured in 
[runner.go](runner.go). 

## AWS 

AWS runners are created in EC2 using the AWS SDK. See the [aws_runners](internal/aws/runners) 
package for details on the implementation. 

## CodeMagic

CodeMagic runners are actually CodeMagic Builds (CI jobs in their service) that then pull the
GitHub Self-Hosted binaries and register themselves as ephemeral (single-use) runners - picking
up a single job from the calling repo and then terminating.

# Testing the application

**Note**: Before running tests, make sure you run the application with environment variable `WEBHOOK_TOKEN=FAKE_TOKEN`.
You will also need to set `GITHUB_TOKEN` to a PAT for the D6E Automaton and `CODEMAGIC_TOKEN` to the PAT for the 
automaton in CodeMagic. These values can all be found in the 
[github-runner-provisioner-secrets.yaml](/keybase/team/datawireio/skunkworks/github-runner-provisioner-secrets.yaml)
file in Keybase - you will need to base64 decode them before use. 

There are several make targets that will send a request to the provisioner on Skunkworks. Target names 
are `test-<RUNNER_TAG>`.

Target takes a `DRY_RUN` variable that makes the request run in dry-run mode. By default, target sets
`DRY_RUN=true`. To override it set `DRY_RUN=false`

Here are example targets to send requests to a locally-running GRP:

```shell
make test-macOS-arm64 HOSTNAME=http://localhost:8080 DRY_RUN=true
```

```shell
make test-ubuntu-arm64 HOSTNAME=http://localhost:8080 DRY_RUN=true
```

**Note**: Be careful when sending requests to production using an HTTP client, since the `dry-run` 
request parameter defaults to true. This is necessary because we have no way to set GitHub to send this 
parameter.

# Env Vars
The runner provisioner requires the following variables to be configured:
- `GITHUB_TOKEN` - a personal access token with admin access to the repo configuring the runners. 
We use the `D6E-Automaton`'s token in production.
- `WEBHOOK_TOKEN` - the secret used to configure the webhook in GitHub. We use the token stored at 
`/Keybase/team/datawireio/infra/github-runner-provisioner-secrets`
- `CODEMAGIC_TOKEN` - the secret used to authenticate to the CodeMagic build API to trigger M1 runners