package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func createMacM1Runner(ctx context.Context, owner string, repo string, dryRun bool) error {
	userData, err := macRunnerUserData(ctx, owner, repo)
	if err != nil {
		return err
	}

	params := ec2.RunInstancesInput{
		MaxCount:                          &macM1Config.instanceCount,
		MinCount:                          &macM1Config.instanceCount,
		DryRun:                            &dryRun,
		ImageId:                           &macM1Config.imageId,
		InstanceInitiatedShutdownBehavior: macM1Config.shutdownBehavior,
		InstanceType:                      macM1Config.instanceType,
		Placement:                         &macM1Config.placement,
		UserData:                          &userData,
	}

	_, err = ec2Client.RunInstances(ctx, &params)
	if err != nil {
		return err
	}

	return nil
}
