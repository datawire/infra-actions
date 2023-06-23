package aws_runners

import "github.com/aws/aws-sdk-go-v2/service/ec2/types"

type runnerConfig struct {
	imageId              string
	hostResourceGroupArn string
	placement            types.Placement
	instanceCount        int32
	shutdownBehavior     types.ShutdownBehavior
	instanceType         types.InstanceType
	keyName              string
}
