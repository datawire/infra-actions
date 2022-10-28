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
      - uses: ./provision-cluster
        with:
          distribution: ${{ matrix.clusters.distribution }}
          version: ${{ matrix.clusters.version }}
          # Tells provision-cluster where to write the kubeconfig file.
          kubeconfig: kubeconfig.yaml
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
      - uses: ./provision-cluster
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

## Custom Mac M1 runners

There are self-hosted Mac M1 (ARM64) runners that can be used in a workflow by setting the`runs-on` option to `macOS-arm64`.
The following example shows how to use this on a matrix.

```yaml
...
jobs:
...
  my_matrix_job:
    strategy:
      matrix:
        runners:
          - ubuntu-latest
          - macOS-arm64
        clusters:
         - distribution: Kubeception
           version: "1.23"
    runs-on: ${{ matrix.runners }}
    steps:
      # The provision-cluster action will automatically register a cleanup hook to remove the
      # cluster it provisions when the job is done.
      - uses: actions/checkout@v3
      - uses: ./provision-cluster
        with:
          distribution: ${{ matrix.clusters.distribution }}
          version: ${{ matrix.clusters.version }}
          # Tells provision-cluster where to write the kubeconfig file.
          kubeconfig: kubeconfig.yaml
      - run: KUBECONFIG=kubeconifig.yaml make tests
...
```

In the future, we may make additional runners available (like Ubuntu ARM64 or Alpine), depending on the needs of the different teams. 

Note that Mac M1 runners are created on demand, and it will take between 30 minutes and up to 3 hours to have one available 
depending on how many workflows are requesting these runners.

Also, there is a limit of Mac M1 10 runners at any point in time. Any build that requests a Mac M1 during this time will fail.  

# Dev loop
 See [DEVELOPING.md](docs/DEVELOPING.md)