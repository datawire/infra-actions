# Cluster Provisioning

## Example Usage

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

    ## For Kubeception klusters

    # A Kubeception secret token
    kubeceptionToken: ...

    ## For GKE clusters:

    # A json encoded string containing GKE credentials:
    gkeCredentials: ...

    # A json encoded string containing additional GKE cluster configuration.
    # Reference the GKE API for more information.
    gkeConfig: ...
```

## References

- [GKE API](https://cloud.google.com/kubernetes-engine/docs/reference/rest)
