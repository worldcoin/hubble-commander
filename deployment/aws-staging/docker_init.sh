#!/bin/bash

# Prepare the ~/.aws directory for the AWS config and credentials
# Remove cron jobs if the directory exists (this is the case only when restarting the container)
AWS_CONFIG_DIR=/root/.aws
if [[ ! -e ${AWS_CONFIG_DIR} ]]; then
    mkdir -p ${AWS_CONFIG_DIR}
else
    crontab -r
fi

# Create the AWS credentials file
echo "[default]
aws_access_key_id=${AWS_ACCESS_KEY_ID}
aws_secret_access_key=${AWS_SECRET_ACCESS_KEY}
" > /root/.aws/credentials

# Create the AWS config file
echo "[default]
region=${AWS_REGION}
output=json
" > /root/.aws/config

# Run the script that initiates backups of the commander state data
bash /root/scripts/start_commander_backups.sh

# Run the cron daemon in the foreground
crond -n
