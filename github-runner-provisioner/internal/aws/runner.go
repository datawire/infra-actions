package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/datawire/infra-actions/github-runner-provisioner/internal/aws/runners"
	"github.com/datawire/infra-actions/github-runner-provisioner/internal/utils"
)

const LabelTag = "label"
const NameTag = "app"
const ownerTag = "owner"
const repoTag = "repo"

const AppName = "github-runner-provisioner"

func RunnerTags(owner string, repo string, runnerLabel string) []types.TagSpecification {
	tags := []types.TagSpecification{
		{
			ResourceType: types.ResourceTypeInstance,
			Tags: []types.Tag{
				{Key: utils.StrPtr(ownerTag), Value: utils.StrPtr(owner)},
				{Key: utils.StrPtr(repoTag), Value: utils.StrPtr(repo)},
				{Key: utils.StrPtr(NameTag), Value: utils.StrPtr(AppName)},
				{Key: utils.StrPtr(LabelTag), Value: &runnerLabel},
			},
		},
	}
	return tags
}

var runnerParams = map[string]func(string, string, string, string, bool) (ec2.RunInstancesInput, error){
	"macos-arm64":  aws_runners.MacM1RunInstancesInput,
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
