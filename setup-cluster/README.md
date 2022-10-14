# Cluster Setup

Use the `setup-cluster` action as described below:

```yaml
      - uses: ./setup-cluster
        with:
          # Tells setup-cluster which kubeconfig file to use.
          kubeconfig: path/to/kubeconfig.yaml

          # Tells setup-cluster where to find manifests. This can be a URL or a path in the
          # filesystem. The manifests can be either raw yaml or kustomize based manifests.
          manifests: url-or-path-to-manifests
```
