package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/datawire/infra-actions/github-runner-provisioner/internal/utils"
	"log"
	"time"
)

type InstanceDetails struct {
	LaunchTime        *time.Time
	InstanceId        *string
	ActionRunnerLabel *string
}

func GetInstances(filter []types.Filter) ([]*InstanceDetails, error) {
	var nextToken *string
	instancesDetails := []*InstanceDetails{}

	for {
		params := ec2.DescribeInstancesInput{
			Filters:   filter,
			NextToken: nextToken,
		}

		describeInstancesOutput, err := Ec2Client.DescribeInstances(context.Background(), &params)
		if err != nil {
			err = fmt.Errorf("error getting EC2 instance information. %v", err)
			return nil, err
		}

		for _, reservation := range describeInstancesOutput.Reservations {
			for _, instance := range reservation.Instances {
				label, err := getActionRunnerLabel(instance)
				if err != nil {
					log.Printf("Error getting runner tag for instance %s: %v\n", *instance.InstanceId, err)
					continue
				}

				instanceDetails := &InstanceDetails{
					LaunchTime:        instance.LaunchTime,
					InstanceId:        instance.InstanceId,
					ActionRunnerLabel: utils.StrPtr(label),
				}

				instancesDetails = append(instancesDetails, instanceDetails)
			}
		}

		if describeInstancesOutput.NextToken == nil {
			break
		}

		nextToken = describeInstancesOutput.NextToken

	}
	return instancesDetails, nil
}

func getActionRunnerLabel(instance types.Instance) (string, interface{}) {
	for _, tag := range instance.Tags {
		if *tag.Key == LabelTag {
			return *tag.Value, nil
		}
	}

	return "", fmt.Errorf("AWS instance %s does not contain tag %s", *instance.InstanceId, LabelTag)
}
