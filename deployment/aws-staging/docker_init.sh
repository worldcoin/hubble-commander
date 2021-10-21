#!/bin/bash

# Prepare directory for the AWS config and credentials
mkdir /root/.aws

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

# Maybe verify that you have access to the aws bucket before doing anything further?

# Run the script that initiates backups of the commander state data
bash /root/scripts/start_commander_backups.sh

# Run the cron daemon in the foreground
crond -n
