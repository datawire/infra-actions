# Github Actions for Test Matrices

This repository hosts github actions that can be used to provision and configure kubernetes clusters. These are intended to facilitate building out a comprehensive [test matrix](../.github/workflows/matrix.yaml) suitable for use in real-world large scale integration and compatibility testing for both Telepresence and Edge Stack.

The [matrix workflow](../.github/workflows/matrix.yaml) illustrates usage of these actions.

## Cluster Provisioning

The [provision-cluster](../provision-cluster/README.md) action can be used to provision different varieties of clusters:

- Kubeception (k3s based)
- GKE

By including this github action in your workflow you can easily run the same test suite against any supported set of clusters:

```yaml
jobs:
  my_matrix_job_gke:
    strategy:
      matrix:
        clusters:
          - version: "1.26"
          - version: "1.27"
          - version: "1.28"
    steps:
      # The provision-cluster action will automatically register a cleanup hook to remove the
      # cluster it provisions when the job is done.
      - uses: datawire/infra-actions/provision-cluster@v0.2.9
        with:
          distribution: GKE
          version: ${{ matrix.clusters.version }}
          # Tells provision-cluster where to write the kubeconfig file.
          kubeconfig: kubeconfig.yaml
          gkeCredentials: ${{ secrets.GOOGLE_APPLICATION_CREDENTIALS }}
          useAuthProvider: "false"
      - run: make tests

  my_matrix_job_kubeception:
    strategy:
      matrix:
        clusters:
          - version: "1.26"
          - version: "1.27"
          - version: "1.28"
    steps:
      # The provision-cluster action will automatically register a cleanup hook to remove the
      # cluster it provisions when the job is done.
      - uses: datawire/infra-actions/provision-cluster@v0.2.9
        with:
          distribution: Kubeception
          version: ${{ matrix.clusters.version }}
          # Tells provision-cluster where to write the kubeconfig file.
          kubeconfig: kubeconfig.yaml
          kubeceptionToken: ${{ secrets.KUBECEPTION_TOKEN }}
      - run: make tests
```

The following inputs apply only to GKE clusters:

- `useAuthProvider`: If set to "true", Authentication is done using an authentication provider, like the [gke-gcloud-auth-plugin](https://cloud.google.com/blog/products/containers-kubernetes/kubectl-auth-changes-in-gke).

The action returns the following outputs:

- `clusterName`: Name of the cluster.
- `projectId`: For GKE, the project ID. Undefined for other cluster providers.
- `location`: For GKE, the cluster location (region or zone). Undefined for other cluster providers.
