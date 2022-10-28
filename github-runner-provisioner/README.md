# Runner Service

This service is based on the [echo
template](https://github.com/datawire/infrastructure/tree/master/echo). Please view the
[README](https://github.com/datawire/infrastructure/tree/master/echo) for details about the dev loop
and how it works.

# Testing the application

Make target `test-github-provisioner` will send a request to the provisioner on Skunkworks. Be carefull when using this 
since it will provision a Mac M1 runner.

To run tests against a local instance of the provisioner use the HOSTNAME parameter like this:

```shell
 make test-github-provisioner HOSTNAME=http://localhost:8080
```