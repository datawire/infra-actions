# Github Actions for Test Matrices

This repository hosts github actions that can be used to provision and configure kubernetes
clusters. These are intended to facilitate building out a comprehensive [test
matrix](../.github/workflows/matrix.yaml) suitable for use in real-world large scale integration and
compatibility testing for both telepresence and edge-stack.

The [matrix workflow](../.github/workflows/matrix.yaml) illustrates an exemplary usage of these
actions.

## Cluster Provisioning

The [provision-cluster](../provision-cluster/README.md) action can be used to provision different
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
           useAuthProvider: "false"
         - distribution: GKE
           version: "1.21"
           useAuthProvider: "true"
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

          useAuthProvider: ${{ matrix.clusters.useAuthProvider }}
      - run: make tests
...
```

The following inputs apply only to GKE clusters:

`useAuthProvider`: If set to "true", Authentication is done using an authentication provider, like the 
[gke-gcloud-auth-plugin](https://cloud.google.com/blog/products/containers-kubernetes/kubectl-auth-changes-in-gke).


The action returns the following outputs:

`clusterName`: Name of the cluster.

`projectId`: For GKE, the project ID. Undefined for other cluster providers.

`location`: For GKE, the cluster location (region or zone). Undefined for other cluster providers.