package aws

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/datawire/infra-actions/github-runner-provisioner/internal/utils"
)

const labelTag = "label"
const nameTag = "app"
const ownerTag = "owner"
const repoTag = "repo"

const appName = "github-runner-provisioner"

func RunnerTags(owner string, repo string, runnerLabel string) []types.TagSpecification {
	tags := []types.TagSpecification{
		{
			ResourceType: types.ResourceTypeInstance,
			Tags: []types.Tag{
				{Key: utils.StrPtr(ownerTag), Value: utils.StrPtr(owner)},
				{Key: utils.StrPtr(repoTag), Value: utils.StrPtr(repo)},
				{Key: utils.StrPtr(nameTag), Value: utils.StrPtr(appName)},
				{Key: utils.StrPtr(labelTag), Value: &runnerLabel},
			},
		},
	}
	return tags
}
