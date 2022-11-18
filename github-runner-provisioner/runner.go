package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/datawire/infra-actions/github-runner-provisioner/internal/aws"
)

var runners = map[string]func(context.Context, string, string, bool) error{
	macM1RunnerLabel:       createMacM1Runner,
	ubuntuArm64RunnerLabel: createUbuntuArm64Runner,
}

func createMacM1Runner(ctx context.Context, owner string, repo string, dryRun bool) error {
	userData, err := macM1RunnerUserData(ctx, owner, repo)
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
		KeyName:                           &macM1Config.keyName,
		Placement:                         &macM1Config.placement,
		UserData:                          &userData,
		TagSpecifications:                 aws.RunnerTags(owner, repo, macM1RunnerLabel),
	}

	_, err = aws.Ec2Client.RunInstances(ctx, &params)
	if err != nil {
		return err
	}

	return nil
}

func createUbuntuArm64Runner(ctx context.Context, owner string, repo string, dryRun bool) error {
	userData, err := ubuntuArm64UserData(ctx, owner, repo)
	if err != nil {
		return err
	}

	params := ec2.RunInstancesInput{
		MaxCount:                          &ubuntuArm64RunnerConfig.instanceCount,
		MinCount:                          &ubuntuArm64RunnerConfig.instanceCount,
		DryRun:                            &dryRun,
		ImageId:                           &ubuntuArm64RunnerConfig.imageId,
		InstanceInitiatedShutdownBehavior: ubuntuArm64RunnerConfig.shutdownBehavior,
		InstanceType:                      ubuntuArm64RunnerConfig.instanceType,
		KeyName:                           &ubuntuArm64RunnerConfig.keyName,
		Placement:                         &ubuntuArm64RunnerConfig.placement,
		UserData:                          &userData,
		TagSpecifications:                 aws.RunnerTags(owner, repo, ubuntuArm64RunnerLabel),
	}

	_, err = aws.Ec2Client.RunInstances(ctx, &params)
	if err != nil {
		return err
	}

	return nil
}
