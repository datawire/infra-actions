# Runner Service

This service is based on the [echo
template](https://github.com/datawire/infrastructure/tree/master/echo). Please view the
[README](https://github.com/datawire/infrastructure/tree/master/echo) for details about the dev loop
and how it works.

# Testing the application

Make target `test-github-provisioner` will send a request to the provisioner on Skunkworks.  

Target takes a `DRY_RUN` variable that makes the request run in dry-run mode. By default, target sets `DRY_RUN=true`. To override it use:

```shell
make test-github-provisioner HOSTNAME=http://localhost:8080 DRY_RUN=false
```

**Note**: Be careful when sending requests to production using a HTTP client, since the `dry-run` request parameter 
defaults to true. This is necessary because we have no way to set GitHub to send this parameter. 

To run tests against a local instance of the provisioner use the HOSTNAME parameter like this:

```shell
 make test-github-provisioner HOSTNAME=http://localhost:8080
```