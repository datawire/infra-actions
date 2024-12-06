package main

import (
	"context"

	"github.com/datawire/infra-actions/github-runner-provisioner/internal/aws"
)

const (
	ubuntuArm64RunnerLabel = "ubuntu-arm64"
)

var runners = map[string]func(context.Context, string, string, bool) error{
	ubuntuArm64RunnerLabel: createUbuntuArm64Runner,
}

func createUbuntuArm64Runner(ctx context.Context, owner string, repo string, dryRun bool) error {
	token, err := getGitHubRunnerToken(ctx, owner, repo, dryRun)
	if err != nil {
		return err
	}

	err = aws.CreateEC2Runner(ctx, owner, repo, token, ubuntuArm64RunnerLabel, dryRun)
	if err != nil {
		return err
	}

	return nil
}
