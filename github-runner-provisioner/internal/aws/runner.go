package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/datawire/infra-actions/github-runner-provisioner/internal/aws/aws_runners"
)

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
		return err
	}

	return nil
}
