## Cluster Setup

Note: This section describes how this feature would work but it's not implemented.

The [setup-cluster](../setup-cluster/README.md) action can be used to configure a cluster with a given
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
