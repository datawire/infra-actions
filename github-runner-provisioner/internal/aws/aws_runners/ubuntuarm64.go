package aws_runners

import (
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

const AmiUbuntuArm64 = "ami-0f69dd1d0d03ad669"
const ubuntuArm64RunnerInstaller = "https://github.com/actions/runner/releases/download/v2.298.2/actions-runner-linux-arm64-2.298.2.tar.gz"
const ubuntuArm64UserDataTemplate = `#!/bin/bash
set -ex

# Download the latest runner package
cat <<EOF > run_agent.sh
set -ex

cd /home/ubuntu

mkdir -p github-agent
cd github-agent

# Download the latest runner package
curl -o github_runner_installer.tar.gz -L '%[1]s'

# Extract the installer
tar xzf ./github_runner_installer.tar.gz

# Configure the agent
./config.sh --url https://github.com/%[2]s/%[3]s --token %[4]s --unattended --labels %[5]s --ephemeral

# Run the agent
timeout 6h ./run.sh
EOF

chmod +x run_agent.sh
sudo su ubuntu - ./run_agent.sh 2>&1 | tee /var/log/github-agent.log

shutdown -h now
`

var ubuntuArm64RunnerConfig = runnerConfig{
	imageId:          AmiUbuntuArm64,
	placement:        types.Placement{},
	instanceCount:    1,
	shutdownBehavior: "terminate",
	instanceType:     "t4g.xlarge",
	keyName:          "m1_mac_runners",
	// keeping default volume size hardcoded to 50gb
	blockDeviceMappings: &[]types.BlockDeviceMapping{
		{
			DeviceName: &[]string{"/dev/sda1"}[0],
			Ebs: &types.EbsBlockDevice{
				VolumeSize: &[]int32{50}[0],
			},
		},
	},
}

func ubuntuArm64UserData(owner string, repo string, token string, label string) (string, error) {
	userData := fmt.Sprintf(ubuntuArm64UserDataTemplate, ubuntuArm64RunnerInstaller, owner, repo, token, label)

	encodedUserData := base64.StdEncoding.EncodeToString([]byte(userData))
	return encodedUserData, nil
}

func UbuntuArm64RunInstancesInput(owner string, repo string, token string, label string, dryRun bool) (ec2.RunInstancesInput, error) {
	userData, err := ubuntuArm64UserData(owner, repo, token, label)
	if err != nil {
		return ec2.RunInstancesInput{}, err
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
		TagSpecifications:                 runnerTags(owner, repo, label),
		BlockDeviceMappings:               *ubuntuArm64RunnerConfig.blockDeviceMappings,
	}

	return params, nil
}
