package main

import (
	"context"
	"github.com/datawire/infra-actions/github-runner-provisioner/internal/aws"
	"github.com/datawire/infra-actions/github-runner-provisioner/internal/codemagic"
)

const macM1RunnerLabel = "macOS-arm64"
const ubuntuArm64RunnerLabel = "ubuntu-arm64"

var runners = map[string]func(context.Context, string, string, bool) error{
	macM1RunnerLabel:       createMacM1Runner,
	ubuntuArm64RunnerLabel: createUbuntuArm64Runner,
}

func createMacM1Runner(ctx context.Context, owner string, repo string, dryRun bool) error {
	token, err := getGitHubRunnerToken(ctx, owner, repo)
	if err != nil {
		return err
	}

	err = codemagic.CreateMacM1Runner(ctx, owner, repo, token, macM1RunnerLabel, dryRun)
	if err != nil {
		return err
	}

	return nil
}

func createUbuntuArm64Runner(ctx context.Context, owner string, repo string, dryRun bool) error {
	token, err := getGitHubRunnerToken(ctx, owner, repo)
	if err != nil {
		return err
	}

	err = aws.CreateEC2Runner(ctx, owner, repo, token, ubuntuArm64RunnerLabel, dryRun)
	if err != nil {
		return err
	}

	return nil
}
