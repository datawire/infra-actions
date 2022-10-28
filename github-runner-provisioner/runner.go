package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func createMacM1Runner(ctx context.Context, owner string, repo string) error {
	client, err := newAwsClient()
	if err != nil {
		return err
	}

	dryRun := true

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

	output, err := client.RunInstances(ctx, &params)
	if err != nil {
		return err
	}

	fmt.Printf("Output: %v\n", output)
	return nil
}
