# Provision Cluster GitHub Action

## Releasing the provision-cluster GitHub Action

GitHub Actions are released by creating a semver tag and pushing it to GitHub. No additional steps are needed.

Once the tag is pushed, then verify the release by using it in the smoke test workflow. Do this by editing `.github/workflows/smoke.yaml`, search for the uses line and update the version to the newly released tag.

```yaml
jobs:
  release_smoke:
    steps:
      - id: provision
        uses: datawire/infra-actions/provision-cluster@vX.Y.Z
```

Pushing the tag should trigger the release smoke test workflow. Verify that this has in fact passed.
