# Cluster Provisioning

Use the `provision-cluster` action as described below:

```yaml
- uses: ./provision-cluster
  with:
    # Tells the action what kind of cluster to create. One of: Kubeception, GKE, EKS, AKS, OpenShift
    distribution: ...
    # Tells the action what version of cluster to create.
    version: 1.27
    # Tells provision-cluster where to write the kubeconfig file.
    kubeconfig: path/to/kubeconfig.yaml

    ## For kubeception klusters

    # A kubeception secret token
    kubeceptionToken: ...

    ## For GKE clusters:

    # A json encoded string containing GKE credentials:
    gkeCredentials: ...
    # A json encoded string containing additional GKE cluster configuration. See GKE Cluster Config Options section for details.
    gkeConfig: ...
```
