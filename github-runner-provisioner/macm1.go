package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type runnerConfig struct {
	imageId              string
	hostResourceGroupArn string
	placement            types.Placement
	instanceCount        int32
	shutdownBehavior     types.ShutdownBehavior
	instanceType         types.InstanceType
}

var macM1HostResourceGroupArn = "arn:aws:resource-groups:us-east-1:914373874199:group/GitHub-Runners"

var macM1Config = runnerConfig{
	imageId: "ami-01b8fcd5770ceb9c1",
	placement: types.Placement{
		HostResourceGroupArn: &macM1HostResourceGroupArn,
	},
	instanceCount:    1,
	shutdownBehavior: "terminate",
	instanceType:     "mac2.metal",
}

func macRunnerUserData(ctx context.Context, owner string, repo string) (string, error) {
	token, err := getGitHubRunnerToken(ctx, owner, repo)
	if err != nil {
		return "", err
	}

	const labels = "macOS-arm64"
	const runnerInstaller = "https://github.com/actions/runner/releases/download/v2.298.2/actions-runner-osx-arm64-2.298.2.tar.gz"

	userData := fmt.Sprintf(`#!/bin/bash
set -ex

cd /Users/ec2-user

# Download the latest runner package
curl -o github_runner_installer.tar.gz -L '%[1]s'

# Extract the installer
sudo su ec2-user -c "tar xzf ./github_runner_installer.tar.gz"
sudo su ec2-user -c "./config.sh --url https://github.com/%[2]s/%[3]s --token %[4]s --unattended --ephemeral --labels %[5]s" 
sudo su ec2-user -c "./run.sh"
shutdown -h now`, runnerInstaller, owner, repo, token, labels)

	encodedUserData := base64.StdEncoding.EncodeToString([]byte(userData))
	return encodedUserData, nil
}
