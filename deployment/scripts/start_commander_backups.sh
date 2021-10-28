#!/bin/bash

# NOTE:
# Remember to use only absolute paths in this script!

# Runs the cron job 5, 25 and 45 minutes past every hour
CRON_SCHEDULE_EXPRESSION="5,25,45 * * * *"

BACKUP_SCRIPT_PATH=""
BACKUPS_DIR_PATH=""
BADGER_DIR_PATH=""
CHAIN_SPEC_DIR_PATH=""
GETH_CHAINDATA_DIR_PATH=""

PIGZ_PATH=""
PG_DUMP_PATH=""
AWS_PATH=""

AWS_S3_BUCKET=""

CRON_SCRIPT_PATH=""

# Initial message for easier logs readability
echo "Starting the docker container for running backups..."

# Schedule backups
BACKUP_COMMAND="bash ${BACKUP_SCRIPT_PATH} ${BACKUPS_DIR_PATH} ${BADGER_DIR_PATH} ${CHAIN_SPEC_DIR_PATH} ${GETH_CHAINDATA_DIR_PATH} ${PIGZ_PATH} ${PG_DUMP_PATH} ${AWS_PATH} ${AWS_S3_BUCKET} > /proc/1/fd/1 2>&1"

# Check the credentials, permission and connection to the AWS S3 bucket
touch ./test_status
aws s3 cp ./test_status s3://"${AWS_S3_BUCKET}" >/dev/null
CONNECTION_STATUS=$?
rm ./test_status

if [[ "${CONNECTION_STATUS}" -eq 0 ]]; then
    "${CRON_SCRIPT_PATH}" "${CRON_SCHEDULE_EXPRESSION}" "${BACKUP_COMMAND}"
else
    echo "The verification process for the AWS S3 bucket failed. Stopping the container..."

    # Kills the process with ID 1 with the SIGTERM signal
    # This results in stopping/killing the container from the inside
    kill -15 1
fi

