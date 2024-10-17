# Cluster Provisioning

## Example Usage

```yaml
- uses: datawire/infra-actions/provision-cluster@v0.4.0
  with:
    distribution: Kubeception
    version: 1.31
    kubeconfig: kubeconfig.yaml

    kubeceptionToken: ${{ secrets.KUBECEPTION_TOKEN }}
```

```yaml
- uses: datawire/infra-actions/provision-cluster@v0.4.0
  with:
    distribution: GKE
    version: 1.31
    kubeconfig: kubeconfig.yaml

    gkeCredentials: ${{ secrets.GOOGLE_APPLICATION_CREDENTIALS }}
    gkeConfig: '{ "initialNodeCount" : 2 }'
```

## References

- [GKE API](https://cloud.google.com/kubernetes-engine/docs/reference/rest)
