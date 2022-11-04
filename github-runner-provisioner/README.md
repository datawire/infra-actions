# Runner Service

This service is based on the [echo
template](https://github.com/datawire/infrastructure/tree/master/echo). Please view the
[README](https://github.com/datawire/infrastructure/tree/master/echo) for details about the dev loop
and how it works.

# Testing the application

**Note**: Before running tests, make sure you run the application with environment variable `WEBHOOK_TOKEN=FAKE_TOKEN`

There are several make targets that will send a request to the provisioner on Skunkworks. Target names 
are `test-<RUNNER_TAG>`.


For example, to run macOS-arm64 tests execute.

```shell
make test-macOS-arm64 HOSTNAME=http://localhost:8080 DRY_RUN=false
```

Target takes a `DRY_RUN` variable that makes the request run in dry-run mode. By default, target sets 
`DRY_RUN=true`. To override it use:

```shell
make test-ubuntu-arm64 HOSTNAME=http://localhost:8080 DRY_RUN=false
```

**Note**: Be careful when sending requests to production using a HTTP client, since the `dry-run` 
request parameter defaults to true. This is necessary because we have no way to set GitHub to send this 
parameter. 

To run tests against a local instance of the provisioner use the HOSTNAME parameter like this:

```shell
 make test-ubuntu-arm64 HOSTNAME=http://localhost:8080
```

# Env Vars
The runner provisioner requires the following variables to be configured:
- `GITHUB_TOKEN` - a personal access token with admin access to the repo configuring the runners. 
We use the `D6E-Automaton`'s token in production.
- `WEBHOOK_TOKEN` - the secret used to configure the webhook in GitHub. We use the token stored at 
`/Keybase/team/datawireio/infra/github-runner-provisioner-secrets`