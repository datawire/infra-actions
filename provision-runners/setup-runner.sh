#!/bin/bash
#set up dependencies
apt update
apt install jq curl unzip -y
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
./aws/install

#acquire runner token
echo "Creating runner token for ${ORGANIZATION}/${GITHUB_REPOSITORY}"
export GITHUB_RUNNER_TOKEN=$(curl \
  -X POST \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer ${INPUT_GITHUB_TOKEN}" \
  https://api.github.com/repos/${ORGANIZATION}/${GITHUB_REPOSITORY}/actions/runners/registration-token \
  | jq -r ".token" )

#create a userdata file to run on cloud-init for the instance
echo "GENERATING USERDATA FOR ${INPUT_INSTANCE_TYPE}"
echo '#!/bin/bash
    cd /home/ec2-user
    # Download the latest runner package
    curl -o github_runner_installer.tar.gz -L '"${INPUT_GITHUB_RUNNER_INSTALLER}"'
    # Extract the installer
    sudo su ec2-user -c "tar xzf ./github_runner_installer.tar.gz"
    sudo su ec2-user -c "./config.sh --url https://github.com/datawire --token '"${GITHUB_RUNNER_TOKEN}"' --unattended --ephemeral"
    sudo su ec2-user -c "./run.sh"
    shutdown -h now' > userdata.sh

#request the instance
echo "CREATING ${INPUT_INSTANCE_TYPE}"
if [-z ${INPUT_HOST_RESOURCE_GROUP_ARN}]; 
  then export PLACEMENT="HostResourceGroupArn=${INPUT_HOST_RESOURCE_GROUP_ARN}";
  else export PLACEMENT="";
  fi
aws ec2 run-instances --image-id $INPUT_AMI_ID --count 1 --instance-type $INPUT_INSTANCE_TYPE \
    --key-name m1_mac_runners --user-data file://userdata.sh --instance-initiated-shutdown-behavior terminate \
    --placement="$PLACEMENT" \ 
    > instance_info.json
cat instance_info.json
export INSTANCE_ID=$(cat instance_info.json | jq -r ".")

#wait for the instance to become available
echo "Waiting for the instance to come online"
aws ec2 wait instance-status-ok --instance-ids $INSTANCE_ID