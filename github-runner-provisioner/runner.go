package main

import (
	"context"
	"fmt"
	config2 "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

func createMacM1Runner(ctx context.Context, owner string, repo string) error {
	cfg, err := config2.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	client := ec2.NewFromConfig(cfg)

	imageId := "ami-01b8fcd5770ceb9c1"
	hostResourceGroupArn := "arn:aws:resource-groups:us-east-1:914373874199:group/GitHub-Runners"
	placement := types.Placement{HostResourceGroupArn: &hostResourceGroupArn}
	dryRun := true

	token, err := getGitHubRunnerToken(ctx, owner, repo)
	if err != nil {
		return err
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

	instanceCount := int32(1)

	params := ec2.RunInstancesInput{
		MaxCount:                          &instanceCount,
		MinCount:                          &instanceCount,
		DryRun:                            &dryRun,
		ImageId:                           &imageId,
		InstanceInitiatedShutdownBehavior: "terminate",
		InstanceType:                      "mac2.metal",
		Placement:                         &placement,
		UserData:                          &userData,
	}

	output, err := client.RunInstances(ctx, &params)
	if err != nil {
		return err
	}

	fmt.Printf("Output: %v\n", output)
	return nil
}

func getGitHubRunnerToken(ctx context.Context, owner string, repo string) (token string, err error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	registrationToken, _, err := client.Actions.CreateRegistrationToken(ctx, owner, repo)
	if err != nil {
		return "", err
	}

	return *registrationToken.Token, nil
}
