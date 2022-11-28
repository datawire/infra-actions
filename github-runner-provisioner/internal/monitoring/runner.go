package monitoring

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/datawire/infra-actions/github-runner-provisioner/internal/aws"
	"github.com/datawire/infra-actions/github-runner-provisioner/internal/utils"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
)

var instanceFilter = []types.Filter{
	{
		Name:   utils.StrPtr("tag:" + aws.NameTag),
		Values: []string{aws.AppName},
	},
	{
		Name:   utils.StrPtr("instance-state-name"),
		Values: []string{string(types.InstanceStateNameRunning)},
	},
}

func UpdateActionRunnersRuntimeMetric() {
	for {
		instancesDetails, err := aws.GetInstances(instanceFilter)
		if err != nil {
			log.Printf("Error getting instance information. %v\n", err)
			break
		}

		for _, instanceDetails := range instancesDetails {
			secondsSinceLaunch := time.Since(*instanceDetails.LaunchTime).Seconds()
			ActionRunnerRuntime.With(prometheus.Labels{"instance_id": *instanceDetails.InstanceId, "label": *instanceDetails.ActionRunnerLabel}).Set(secondsSinceLaunch)
		}

		time.Sleep(15 * time.Second)
	}
}