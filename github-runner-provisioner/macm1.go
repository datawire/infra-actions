package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

const runnerLabel = "macOS-arm64"
const runnerInstaller = "https://github.com/actions/runner/releases/download/v2.298.2/actions-runner-osx-arm64-2.298.2.tar.gz"
const userDataTemplate = `#!/bin/bash
set -ex

cd /Users/ec2-user

brew install coreutils

# Download the latest runner package
curl -o github_runner_installer.tar.gz -L '%[1]s'

# Extract the installer
sudo su ec2-user -c "tar xzf ./github_runner_installer.tar.gz"
sudo su ec2-user -c "./config.sh --url https://github.com/%[2]s/%[3]s --token %[4]s --unattended --labels %[5]s" 

# Run agent for to 24 hours
sudo su ec2-user -c "timeout 24h ./run.sh"

shutdown -h now`

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
	imageId: AMI_MACOS_12_6_ARM64,
	placement: types.Placement{
		HostResourceGroupArn: &macM1HostResourceGroupArn,
	},
	instanceCount:    1,
	shutdownBehavior: "terminate",
	instanceType:     "mac2.metal",
	keyName:          "m1_mac_runners",
}

func macRunnerUserData(ctx context.Context, owner string, repo string) (string, error) {
	token, err := getGitHubRunnerToken(ctx, owner, repo)
	if err != nil {
		return "", err
	}

	userData := fmt.Sprintf(userDataTemplate, runnerInstaller, owner, repo, token, runnerLabel)

	encodedUserData := base64.StdEncoding.EncodeToString([]byte(userData))
	return encodedUserData, nil
}
