#!/bin/bash
echo "gh repo token ${GITHUB_REPO_TOKEN}"
echo "gh repo ${GITHUB_REPOSITORY}"
set -e
exit
#set up dependencies

apt -qq update
apt -qq install jq curl unzip -y
# curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
# unzip awscliv2.zip
# ./aws/install

#acquire runner token
echo "Creating runner token for ${GITHUB_REPOSITORY}"
curl -s -X POST \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer ${GITHUB_REPO_TOKEN}" \
  https://api.github.com/repos/${GITHUB_REPOSITORY}/actions/runners/registration-token > token_output.json
cat token_output.json
export GITHUB_RUNNER_TOKEN=$(cat token_output.json | jq -r ".token")

#create a userdata file to run on cloud-init for the instance
echo "GENERATING USERDATA FOR ${INSTANCE_TYPE}"
echo '#!/bin/bash
    cd /home/ec2-user
    # Download the latest runner package
    curl -o github_runner_installer.tar.gz -L '"${GITHUB_RUNNER_INSTALLER}"'
    # Extract the installer
    sudo su ec2-user -c "tar xzf ./github_runner_installer.tar.gz"
    sudo su ec2-user -c "./config.sh --url https://github.com/datawire --token '"${GITHUB_RUNNER_TOKEN}"' --unattended --ephemeral"
    sudo su ec2-user -c "./run.sh"
    shutdown -h now' > userdata.sh

cat userdata.sh
exit

#request the instance
echo "CREATING ${INSTANCE_TYPE}"
if [ -z "${HOST_RESOURCE_GROUP_ARN}" ]; 
  then export PLACEMENT="";
  else export PLACEMENT="--placement=\"HostResourceGroupArn=${HOST_RESOURCE_GROUP_ARN}\"";
  fi
aws ec2 run-instances --image-id $AMI_ID --count 1 --instance-type $INSTANCE_TYPE \
    --key-name m1_mac_runners --user-data file://userdata.sh --instance-initiated-shutdown-behavior terminate \
    $PLACEMENT \ 
    > instance_info.json
cat instance_info.json
export INSTANCE_ID=$(cat instance_info.json | jq -r ".")

#wait for the instance to become available
echo "Waiting for the instance to come online"
aws ec2 wait instance-status-ok --instance-ids $INSTANCE_ID