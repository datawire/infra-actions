# Custom GitHub action runners

There are Mac M1 and Ubuntu ARM64 runners available for GitHub actions.

In the future, we may make additional runners available depending on the needs of the different teams.

## Mac M1 runners

There are self-hosted Mac M1 (ARM64) runners that can be used in a workflow by using `runs-on: macOS-arm64`.

The following example shows how to use this on a matrix.

```yaml
...
jobs:
    my_job:
      runs-on: macOS-arm64
      steps:
        # The provision-cluster action will automatically register a cleanup hook to remove the
        # cluster it provisions when the job is done.
        - uses: actions/checkout@v3
        ...
```

Note that Mac M1 runners are created on demand, and it will take between 30 minutes and up to 3 hours to have one available
depending on how many jobs are running on them.

Also, there is a limit of 10 Mac M1 runners at any point in time. Any build that requests a Mac M1 during this time will stay in a queued state and may fail.

## Ubuntu ARM64 runners

These self-hosted runners run on AWS and are created on-demand. It takes about a minute for the runner to be available, and once the job finishes, it's destroyed.

To request one, use label `ubuntu-arm64`:

```yaml
...
jobs:
  my_job:
    runs-on: ubuntu-arm64
    steps:
      # The provision-cluster action will automatically register a cleanup hook to remove the
      # cluster it provisions when the job is done.
      - uses: actions/checkout@v3
      ...
```