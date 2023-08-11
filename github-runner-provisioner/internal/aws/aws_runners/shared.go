package aws_runners

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/datawire/infra-actions/github-runner-provisioner/internal/utils"
)

type runnerConfig struct {
	imageId              string
	hostResourceGroupArn string
	placement            types.Placement
	instanceCount        int32
	shutdownBehavior     types.ShutdownBehavior
	instanceType         types.InstanceType
	keyName              string
	blockDeviceMappings  *[]types.BlockDeviceMapping
}

const LabelTag = "label"
const NameTag = "app"
const ownerTag = "owner"
const repoTag = "repo"

const AppName = "github-runner-provisioner"

func runnerTags(owner string, repo string, runnerLabel string) []types.TagSpecification {
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
