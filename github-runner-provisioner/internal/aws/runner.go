package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/smithy-go"
	"github.com/datawire/infra-actions/github-runner-provisioner/internal/aws/aws_runners"
)

const dryRunApiError = "DryRunOperation"

var runnerParams = map[string]func(string, string, string, string, bool) (ec2.RunInstancesInput, error){
	"macOS-arm64":  aws_runners.MacM1RunInstancesInput,
	"ubuntu-arm64": aws_runners.UbuntuArm64RunInstancesInput,
}

func CreateEC2Runner(ctx context.Context, owner string, repo string, token string, label string, dryRun bool) error {
	params, err := runnerParams[label](owner, repo, token, label, dryRun)
	if err != nil {
		return err
	}

	ec2Client := NewEc2Client()

	_, err = ec2Client.Client.RunInstances(ctx, &params)
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			code := apiErr.ErrorCode()
			if code == dryRunApiError && dryRun {
				return nil
			}
		}
	}
	return nil
}
