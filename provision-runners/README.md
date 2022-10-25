# Provision Runners

Use the `provision-runners` action in the following way: 
```yaml
        - uses: ./provision-runners
          with: 
            aws_access_key_id: #'Use a secret - this account needs "ec2-admin" permissions'
            aws_secret_access_key: #'Use a secret - this account needs "ec2-admin" permissions'
            github_repo_token: #"Use a secret. Token used to register a github action runner. If D6E Automaton token make sure it is added to the repo with Admin powers."
            region: #'aws region for the runner. Default us-east-1'
            ami_id: #"ami ID that will run on the mac2 instance. Must be associated with the Host Resource Group license."
            instance_type: #"the instance type to spin up in the Host Resource Group"
            host_resource_group_arn: #"the arn of the HRG that will automate setup of the dedicated hosts for the runners"
                #default: "arn:aws:resource-groups:us-east-1:914373874199:group/GitHub-Runners"
            github_runner_installer: #"the url of the tar archive for the github runner installer. Must match desired instance architecture"
                #default: https://github.com/actions/runner/releases/download/v2.298.2/actions-runner-osx-arm64-2.298.2.tar.gz
  
```