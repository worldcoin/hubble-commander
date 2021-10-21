#!/bin/bash

mkdir /root/.aws

echo "[default]
aws_access_key_id=${AWS_ACCESS_KEY_ID}
aws_secret_access_key=${AWS_SECRET_ACCESS_KEY}
" > /root/.aws/credentials

echo "[default]
region=${AWS_REGION}
output=json
" > /root/.aws/config

# Maybe verify that you have access to the aws bucket before doing anything further?

bash /root/scripts/start_commander_backups.sh

crond -n
