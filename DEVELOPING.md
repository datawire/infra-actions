# Documentation to enable developing and releasing the items in this repository.

## Releasing the provision-cluster github action:

Github actions are released by creating a semver tag and pusing it to github. No additional steps
are needed.

### Step 1: Query existing tags

Use `git tag -l` to find existing tag names. Releases are 

### Step 2: Tag with your new version number

Use `git tag vX.Y.Z` to tag with your new version number, and then run `git push --tags` to push the
new tag up to github.

### Step 3: Verify the release works by updating a workflow

Pushing the tag should trigger the release smoke test workflow. Verify that this has in fact passed.
