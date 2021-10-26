#!/bin/bash

# NOTE:
# Remember to use only absolute paths in this script!

# Runs the cron job every 20 minutes
CRON_SCHEDULE_EXPRESSION="*/20 * * * *"

BACKUP_SCRIPT_PATH=""
COMMANDER_DIR_PATH=""
BADGER_DIR_PATH=""
CHAIN_SPEC_DIR_PATH=""
GETH_CHAINDATA_DIR_PATH=""

PIGZ_PATH=""
PG_DUMP_PATH=""

AWS_S3_BUCKET=""

CRON_SCRIPT_PATH=""

BACKUP_COMMAND="bash ${BACKUP_SCRIPT_PATH} ${COMMANDER_DIR_PATH} ${BADGER_DIR_PATH} ${CHAIN_SPEC_DIR_PATH} ${GETH_CHAINDATA_DIR_PATH} ${PIGZ_PATH} ${PG_DUMP_PATH} ${AWS_S3_BUCKET}"

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

