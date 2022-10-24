#!/bin/bash
cd /home/ec2-user
# Download the latest runner package
curl -O -L https://github.com/actions/runner/releases/download/v2.298.2/actions-runner-linux-x64-2.298.2.tar.gz
# Extract the installer
sudo su ec2-user -c "tar xzf ./actions-runner-linux-x64-2.298.2.tar.gz"
sudo su ec2-user -c "./config.sh --url https://github.com/datawire --token GITHUB_RUNNER_TOKEN --unattended --ephemeral"
sudo su ec2-user -c "./run.sh"
shutdown -h now