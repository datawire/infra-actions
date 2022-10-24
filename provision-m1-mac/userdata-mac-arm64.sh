#!/bin/zsh
cd /Users/ec2-user
# Download the latest runner package
curl -o actions-runner-osx-arm64-2.298.2.tar.gz -L https://github.com/actions/runner/releases/download/v2.298.2/actions-runner-osx-arm64-2.298.2.tar.gz
echo "e124418a44139b4b80a7b732cfbaee7ef5d2f10eab6bcb3fd67d5541493aa971  actions-runner-osx-arm64-2.298.2.tar.gz" | shasum -a 256 -c
sudo su ec2-user -c "tar xzf ./actions-runner-osx-arm64-2.298.2.tar.gz"
sudo su ec2-user -c "./config.sh --url https://github.com/datawire --token GITHUB_RUNNER_TOKEN --unattended --ephemeral"
# sudo su ec2-user -c "./run.sh"
# shutdown -h now