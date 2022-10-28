package main

import (
	"context"
	"encoding/base64"
	"fmt"
	config2 "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
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

var amcM1HostResourceGroupArn = "arn:aws:resource-groups:us-east-1:914373874199:group/GitHub-Runners"

var macM1Config = runnerConfig{
	imageId: "ami-01b8fcd5770ceb9c1",
	placement: types.Placement{
		HostResourceGroupArn: &amcM1HostResourceGroupArn,
	},
	instanceCount:    1,
	shutdownBehavior: "terminate",
	instanceType:     "mac2.metal",
}

func newAwsClient() (*ec2.Client, error) {
	cfg, err := config2.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := ec2.NewFromConfig(cfg)
	return client, nil
}

func macRunnerUserData(ctx context.Context, owner string, repo string) (string, error) {
	token, err := getGitHubRunnerToken(ctx, owner, repo)
	if err != nil {
		return "", err
	}

	runnerInstaller := "https://github.com/actions/runner/releases/download/v2.298.2/actions-runner-osx-arm64-2.298.2.tar.gz"
	userData := fmt.Sprintf(`#!/bin/bash
    cd '"${AMI_HOME}"'/'"${AMI_USER}"'
    # Download the latest runner package
    curl -o github_runner_installer.tar.gz -L '"%s"'
    # Extract the installer
    sudo su '"${AMI_USER}"' -c "tar xzf ./github_runner_installer.tar.gz"
    sudo su '"${AMI_USER}"' -c "./config.sh --url https://github.com/%s --token '"%s"' --unattended --ephemeral --labels '"${RUNNER_LABELS}"'" 
    sudo su '"${AMI_USER}"' -c "./run.sh"
    shutdown -h now`, repo, runnerInstaller, token)
	encodedUserData := base64.StdEncoding.EncodeToString([]byte(userData))
	return encodedUserData, nil
}
