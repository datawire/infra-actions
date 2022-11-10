# Custom GitHub action runners

There are self-hosted Mac M1 and Ubuntu ARM64 runners available for GitHub actions. There runners are EC2 instances 
hosted in AWS.

In the future, we may make additional runners available depending on the needs of the different teams.

## Repository configuration

Before a job can use a self-hosted runner, the following settings need to be configured in the GitHub repository:
 
1. Add the `d6e-automaton` account as a repo administrator (`Repo -> Settings -> Collaborators and teams`)
2. Add a webhook (`Repo -> Settings -> Webhooks`) with the following settings:
   1. Payload URL: https://sw.bakerstreet.io/github-runner-provisioner/
   2. Content type: application/x-www-form-urlencoded
   3. Secret: Enter the value found in `/keybase/team/datawireio/secrets/github-actions/github-infra-actions`
   4. SSL verification: Enable
   5. Which events trigger the webhook? -> Let me select individual events ->  Workflow jobs

Once the webhook is configured, you can use the runners as described below.

## Mac M1 runners

There are self-hosted Mac M1 (ARM64) runners that can be used in a workflow by using `runs-on: macOS-arm64`.

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

The following limitations apply to Mac M1 runners:
- It will take between 30 minutes and up to 3 hours for a runner to be available from the moment it's requested by a job.
- There is a limit of 10 Mac M1 runners at any point in time. Any build that requests a Mac M1 during this time will 
  stay in a queued state until aruner is available. If a job is queued for more than 24 hours, it will be marked as failed.
- Once a Mac M1 runner is created, it will continue to run for up to 24 hours, picking-up oe or more jobs. What the means 
  is that jobs are responsible for ensuring that runners are in a clean state before they are used.

## Ubuntu ARM64 runners

These self-hosted runners are created on-demand. It takes about a minute for the runner to be available, and once the 
job finishes, they are destroyed.

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