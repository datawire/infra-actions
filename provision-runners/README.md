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

# Implementation details

We run our hosted runners in AWS. Our currently supported runner types are `mac2.metal` and `ubuntu:22.04` running on arm64/x86_64 instances. 

For any instance the action will need an AMI to launch that instance with, as well as secrets configured for AWS authentication and a GitHub 
token with admin priveleges for the repository that is having runners added. Details about that AMI are also important - the home directory and
user are needed to get the github runner and and going. Allocated instances will automatically clean themselves up after they handle a job.

## M1 Macs

To scale up M1 mac instances we have to use `dedicated hosts` on AWS to reserve a physical mac computer, which we can then schedule instances 
on top of. We automate the host management using `host resource groups` - a system that will schedule a dedicated host when an instance 
matching the HRG's license is requested. That license is tied to an `ami` that we have to own. NOTE: When an M1 instance is requested the host
will be allocated for a minimum of 24 hours - additional jobs running on that host will have no added cost, but will have a 1-3 hour delay 
between instances for mandatory host cleaning. An instance must be allocated for every job. 

Below are the AMIs for our ARM mac instances: 
- macOS 12.6 on arm64: ami-01b8fcd5770ceb9c1

## Ubuntu images

Ubuntu instances are allocated as regular EC2 instances. They spin up in < 3min and clean up equally fast. 

Below are the AMIs that have been tested for Ubuntu instances: 
- ubuntu 22.04 on arm64: ami-0f69dd1d0d03ad669 (tested on t4g.micro instances)