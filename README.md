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

# Dev loop
 See [DEVELOPING.md](docs/DEVELOPING.md)