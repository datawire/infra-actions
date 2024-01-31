# Working with GitHub workflows and actions

GitHub workflows and any actions used by them can be tested locally using [act](https://github.com/nektos/act).

Once `act` is installed, it can be invoked from the repository root like this:

```shell
act pull_request
```

`act` can pass secrets with the command line option `-s`. For example, to pass a secret called `KUBECEPTION_TOKEN` run it like this:

```shell
act pull_request -s KUBECEPTION_TOKEN=MY_TOKEN
```
