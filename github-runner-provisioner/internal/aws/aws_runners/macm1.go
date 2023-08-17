package aws_runners

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

const AmiMacOs12_6Arm64 = "ami-01b8fcd5770ceb9c1"
const macM1RunnerInstaller = "https://github.com/actions/runner/releases/download/v2.308.0/actions-runner-osx-x64-2.308.0.tar.gz"
const macM1UserDataTemplate = `#!/bin/bash
set -x

cd /Users/ec2-user

# Download the latest runner package
cat <<EOF > run_agent.sh
set -ex

mkdir -p github-agent
cd github-agent

/opt/homebrew/bin/brew install coreutils

# Download the latest runner package
curl -o github_runner_installer.tar.gz -L '%[1]s'

# Extract the installer
tar xzf ./github_runner_installer.tar.gz

# Configure the agent
./config.sh --url https://github.com/%[2]s/%[3]s --token %[4]s --unattended --labels %[5]s 

# Run the agent for 1 day
/opt/homebrew/bin/timeout 1d ./run.sh || true

# De-register the agent
./config.sh remove --token %[4]s
EOF

chmod +x run_agent.sh
sudo su ec2-user - ./run_agent.sh 2>&1 | tee /var/log/github-agent.log
shutdown -h now
`

var macM1HostResourceGroupArn = "arn:aws:resource-groups:us-east-1:914373874199:group/GitHub-Runners"
var macM1AvailabilityZone = "us-east-1a"

var macM1RunnerConfig = runnerConfig{
	imageId: AmiMacOs12_6Arm64,
	placement: types.Placement{
		HostResourceGroupArn: &macM1HostResourceGroupArn,
		AvailabilityZone:     &macM1AvailabilityZone,
	},
	instanceCount:    1,
	shutdownBehavior: "terminate",
	instanceType:     "mac2.metal",
	keyName:          "m1_mac_runners",
}

func macM1RunnerUserData(owner string, repo string, token string, label string) (string, error) {
	userData := fmt.Sprintf(macM1UserDataTemplate, macM1RunnerInstaller, owner, repo, token, label)

	encodedUserData := base64.StdEncoding.EncodeToString([]byte(userData))
	return encodedUserData, nil
}

func MacM1RunInstancesInput(owner string, repo string, token string, label string, dryRun bool) (ec2.RunInstancesInput, error) {
	userData, err := macM1RunnerUserData(owner, repo, token, label)
	if err != nil {
		return ec2.RunInstancesInput{}, err
	}

	params := ec2.RunInstancesInput{
		MaxCount:                          &macM1RunnerConfig.instanceCount,
		MinCount:                          &macM1RunnerConfig.instanceCount,
		DryRun:                            &dryRun,
		ImageId:                           &macM1RunnerConfig.imageId,
		InstanceInitiatedShutdownBehavior: macM1RunnerConfig.shutdownBehavior,
		InstanceType:                      macM1RunnerConfig.instanceType,
		KeyName:                           &macM1RunnerConfig.keyName,
		Placement:                         &macM1RunnerConfig.placement,
		UserData:                          &userData,
		TagSpecifications:                 runnerTags(owner, repo, label),
	}

	return params, nil
}
