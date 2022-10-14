# Cluster Provisioning

Use the `provision-cluster` action as described below:

```yaml
      - uses: ./provision-cluster
        with:
          # Tells the action what kind of cluster to create. One of: Kubeception, GKE, EKS, AKS, OpenShift
          distribution: ...
          # Tells the action what version of cluster to create.
          version: 1.23
          # Tells provision-cluster where to write the kubeconfig file.
          kubeconfig: path/to/kubeconfig.yaml
```
