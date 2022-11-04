package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

const AmiMacOs12_6Arm64 = "ami-01b8fcd5770ceb9c1"
const macM1RunnerLabel = "macOS-arm64"
const macM1RunnerInstaller = "https://github.com/actions/runner/releases/download/v2.298.2/actions-runner-osx-arm64-2.298.2.tar.gz"
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

type runnerConfig struct {
	imageId              string
	hostResourceGroupArn string
	placement            types.Placement
	instanceCount        int32
	shutdownBehavior     types.ShutdownBehavior
	instanceType         types.InstanceType
	keyName              string
}

var macM1HostResourceGroupArn = "arn:aws:resource-groups:us-east-1:914373874199:group/GitHub-Runners"

var macM1Config = runnerConfig{
	imageId: AmiMacOs12_6Arm64,
	placement: types.Placement{
		HostResourceGroupArn: &macM1HostResourceGroupArn,
	},
	instanceCount:    1,
	shutdownBehavior: "terminate",
	instanceType:     "mac2.metal",
	keyName:          "m1_mac_runners",
}

func macM1RunnerUserData(ctx context.Context, owner string, repo string) (string, error) {
	token, err := getGitHubRunnerToken(ctx, owner, repo)
	if err != nil {
		return "", err
	}

	userData := fmt.Sprintf(macM1UserDataTemplate, macM1RunnerInstaller, owner, repo, token, macM1RunnerLabel)

	encodedUserData := base64.StdEncoding.EncodeToString([]byte(userData))
	return encodedUserData, nil
}
