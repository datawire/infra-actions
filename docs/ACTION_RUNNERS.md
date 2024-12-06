# Custom GitHub action runners

Self-hosted Ubuntu ARM64 runners are available for GitHub actions. These runners are EC2 instances hosted in AWS.

In the future, we may make additional runners available depending on the needs of the different teams.

## Repository configuration

Before a job can use a self-hosted runner, the following settings need to be configured in the GitHub repository:

1. Add the `d6e-automaton` account as a repo administrator (`Repo ⇾ Settings ⇾ Collaborators and teams`)
2. Add a webhook (`Repo ⇾ Settings ⇾ Webhooks`) with the following settings:
   1. Payload URL: `https://sw.bakerstreet.io/github-runner-provisioner/`
   2. Content type: `application/x-www-form-urlencoded`
   3. Secret: Enter the value found in `/keybase/team/datawireio/secrets/github-actions/github-infra-actions`
   4. SSL verification: `Enable`
   5. Which events trigger the webhook? `Let me select individual events` ⇾ `Workflow jobs`

Once the webhook is configured, you can use the runners as described below.

## Ubuntu ARM64 runners

These self-hosted runners are created on-demand. It takes about a minute for the runner to be available, and once the job finishes, they are destroyed.

To request one, use label `ubuntu-arm64`:

```yaml
jobs:
  my_job:
    runs-on: ubuntu-arm64
    steps:
      - uses: actions/checkout@v4
```
