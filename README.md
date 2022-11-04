# Github Actions for Test Matrices

This repository hosts github actions that can be used to provision and configure kubernetes
clusters. These are intended to facilitate building out a comprehensive [test
matrix](.github/workflows/matrix.yaml) suitable for use in real-world large scale integration and
compatibility testing for both telepresence and edge-stack.

The [matrix workflow](.github/workflows/matrix.yaml) illustrates an exemplary usage of these
actions.

## Cluster Provisioning

The [provision-cluster](provision-cluster/README.md) action can be used to provision different
varieties of clusters:

 - Kubeception (k3s based)
 - GKE
 - EKS (unimplemented)
 - AKS (unimplemented)

By including this github action in your workflow you can easily run the same test suite against any
supported set of clusters:

```yaml
...
jobs:
...
  my_matrix_job:
    strategy:
      matrix:
        clusters:
         - distribution: GKE
           version: "1.21"
         - distribution: AKS
           version: "1.22"
         - distribution: Kubeception
           version: "1.23"
    steps:
      # The provision-cluster action will automatically register a cleanup hook to remove the
      # cluster it provisions when the job is done.
      - uses: datawire/infra-actions/provision-cluster@v0.2.0
        with:
          distribution: ${{ matrix.clusters.distribution }}
          version: ${{ matrix.clusters.version }}
          # Tells provision-cluster where to write the kubeconfig file.
          kubeconfig: kubeconfig.yaml

          kubeceptionToken: ${{ secrets.KUBECEPTION_TOKEN }}
          gkeCredentials: ${{ secrets.GOOGLE_APPLICATION_CREDENTIALS }}
      - run: KUBECONFIG=kubeconifig.yaml make tests
...
```

## Cluster Setup

XXX: unimplemented

The [setup-cluster](setup-cluster/README.md) action can be used to configure a cluster with a given
set of manifests required for test execution. The action will not only intelligently apply the
manifests (dealing with any interdependencies), but also ensure that all deployments, statefulsets,
daemonsets, etc, are fully available, ready, and passing their health checks before allowing the job
to proceed:

```yaml
...
jobs:
...
  my_matrix_job:
    ...
    steps:
      - uses: datawire/infra-actions/provision-cluster@v0.2.0
        with:
          distribution: ${{ matrix.clusters.distribution }}
          version: ${{ matrix.clusters.version }}
          # Tells provision-cluster where to write the kubeconfig file.
          kubeconfig: kubeconfig.yaml
          # For convenience, the provision-cluster manifest will invoke the setup-cluster manifest if you
          # pass in a pointer to manifests, however you can also use it independently as shown below just
          # in case your cluster does not come from the ./provision-cluster action.
          #
          # manifests: <url-or-file>
      - uses: ./setup-cluster
        with:
          # Tells setup-cluster which kubeconfig file to use.
          kubeconfig: kubeconfig.yaml
          # The manifests parameter can point to a file or url and can include raw yaml or kustomized manifests.
          manifests: https://github.com/datawire/my-repo/manifests.yaml
      - run: KUBECONFIG=kubeconifig.yaml make tests
...
```

# Custom GitHub action runners

There are Mac M1 and Ubuntu ARM64 runners available for GitHub actions.

In the future, we may make additional runners available depending on the needs of the different teams.

## Mac M1 runners

There are self-hosted Mac M1 (ARM64) runners that can be used in a workflow by using `runs-on: macOS-arm64`.

The following example shows how to use this on a matrix.

```yaml
...
jobs:
    my_job:
      runs-on: macOS-arm64
      steps:
        # The provision-cluster action will automatically register a cleanup hook to remove the
        # cluster it provisions when the job is done.
        - uses: actions/checkout@v3
        ...
```

Note that Mac M1 runners are created on demand, and it will take between 30 minutes and up to 3 hours to have one available 
depending on how many jobs are running on them.

Also, there is a limit of 10 Mac M1 runners at any point in time. Any build that requests a Mac M1 during this time will stay in a queued state and may fail.

## Ubuntu ARM64 runners

These self-hosted runners run on AWS and are created on-demand. It takes about a minute for the runner to be available, and once the job finishes, it's destroyed.

To request one, use label `ubuntu-arm64`:

```yaml
...
jobs:
  my_job:
    runs-on: ubuntu-arm64
    steps:
      # The provision-cluster action will automatically register a cleanup hook to remove the
      # cluster it provisions when the job is done.
      - uses: actions/checkout@v3
      ...
```

# Dev loop
 See [DEVELOPING.md](docs/DEVELOPING.md)
