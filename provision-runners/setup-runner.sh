#!/bin/bash
[[ ! -z "$RUNNER_ADMIN_TOKEN" ]] && echo "RUNNER_ADMIN_TOKEN Not empty" || echo "RUNNER_ADMIN_TOKEN Empty"
[[ ! -z "$AMI_ID" ]] && echo "AMI_ID Not empty" || echo "AMI_ID Empty"
[[ ! -z "$AWS_ACCESS_KEY_ID" ]] && echo "AWS_ACCESS_KEY_ID Not empty" || echo "AWS_ACCESS_KEY_ID Empty"
[[ ! -z "$AWS_SECRET_ACCESS_KEY" ]] && echo "AWS_SECRET_ACCESS_KEY Not empty" || echo "AWS_SECRET_ACCESS_KEY Empty"
[[ ! -z "$AWS_DEFAULT_REGION" ]] && echo "AWS_DEFAULT_REGION Not empty" || echo "AWS_DEFAULT_REGION Empty"

set -e

echo "setting up dependencies"
sudo apt -qq update
sudo apt -qq install jq curl unzip moreutils -y &> /dev/null

#Only install awscli if it's missing (running act)
if [ ! -x "$(command -v aws)" ]
  then 
    echo "aws-cli was not found - installing now"
    curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
    unzip -qq awscliv2.zip
    ./aws/install
  else echo "aws-cli was found on system"
fi 

#acquire runner token
echo "Creating runner token for ${GITHUB_REPOSITORY}"
curl -s -X POST \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer ${RUNNER_ADMIN_TOKEN}" \
  https://api.github.com/repos/${GITHUB_REPOSITORY}/actions/runners/registration-token > token_output.json

export GITHUB_RUNNER_TOKEN=$(cat token_output.json | jq -r ".token")
if [ ${GITHUB_RUNNER_TOKEN} == "null" ]; 
  then echo "Could not generate valid token. Make sure the token has admin access to the repository." && exit 1; 
  else echo "Successfully generated runner token.";
fi

#create a userdata file to run on cloud-init for the instance
echo "GENERATING USERDATA FOR ${INSTANCE_TYPE}"
echo '#!/bin/bash
    cd '"${AMI_HOME}"'/'"${AMI_USER}"'
    # Download the latest runner package
    curl -o github_runner_installer.tar.gz -L '"${GITHUB_RUNNER_INSTALLER}"'
    # Extract the installer
    sudo su '"${AMI_USER}"' -c "tar xzf ./github_runner_installer.tar.gz"
    sudo su '"${AMI_USER}"' -c "./config.sh --url https://github.com/'"${GITHUB_REPOSITORY}"' --token '"${GITHUB_RUNNER_TOKEN}"' --unattended --ephemeral --labels '"${RUNNER_LABELS}"'" 
    sudo su '"${AMI_USER}"' -c "./run.sh"
    shutdown -h now' > userdata.sh

#request the instance
echo "CREATING ${INSTANCE_TYPE}"
if [ -z "${HOST_RESOURCE_GROUP_ARN}" ]; 
  then export PLACEMENT="" && echo "No HRG given - will create floating instance";
  else export PLACEMENT="--placement=HostResourceGroupArn=${HOST_RESOURCE_GROUP_ARN}";
fi

aws ec2 run-instances --image-id $AMI_ID --count 1 --instance-type $INSTANCE_TYPE \
  --key-name m1_mac_runners --user-data file://userdata.sh \
  --instance-initiated-shutdown-behavior terminate $PLACEMENT > instance_info.json
cat instance_info.json
export INSTANCE_ID=$(cat instance_info.json | jq -r ".Instances[0].InstanceId")

echo "Instance ${INSTANCE_ID} will be ready to run jobs shortly"