package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

const AmiUbuntuArm64 = "ami-0f69dd1d0d03ad669"
const ubuntuArm64RunnerLabel = "ubuntu-arm64"
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
	instanceType:     "t4g.medium",
	keyName:          "m1_mac_runners",
}

func ubuntuArm64UserData(ctx context.Context, owner string, repo string) (string, error) {
	token, err := getGitHubRunnerToken(ctx, owner, repo)
	if err != nil {
		return "", err
	}

	userData := fmt.Sprintf(ubuntuArm64UserDataTemplate, ubuntuArm64RunnerInstaller, owner, repo, token, ubuntuArm64RunnerLabel)

	encodedUserData := base64.StdEncoding.EncodeToString([]byte(userData))
	return encodedUserData, nil
}
