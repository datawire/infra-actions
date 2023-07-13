# Documentation to enable developing and releasing the items in this repository.

## Releasing the provision-cluster GitHub Action:

GitHub Actions are released by creating a semver tag and pushing it to GitHub. No additional steps
are needed.

### Step 1: Query existing tags

Use `git pull` to make sure you have all tags locally and then use `git tag -l` to find existing tag
names. Release tags are of the form `vX.Y.Z` and release versions should follow semver.

### Step 2: Tag with your new version number

Use `git tag vX.Y.Z` to tag with your new version number, and then run `git push --tags` to push the
new tag up to GitHub.

### Step 3: Verify the release works by updating the smoke test workflow.

Once the tag is pushed, then verify the release by using it in the smoke test workflow. Do this by
editing `.github/workflows/smoke.yaml`, search for the uses line and update the version to the newly
released tag, e.g.:

```
...
       - uses: datawire/infra-actions/provision-cluster@vX.Y.Z
...
```

Pushing the tag should trigger the release smoke test workflow. Verify that this has in fact passed.
