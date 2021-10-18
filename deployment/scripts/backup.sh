#!/bin/bash

if [[ "$#" -ne 6 && "$#" -ne 7 ]]; then
    echo "Script used for backing up the current state of the commander."
    echo "This script generates a .tgz file with a compressed backup."
    echo "Optionally, it can upload said backup to an AWS S3 bucket."
    echo "Use the unpigz.sh script to decompress a backup file."
    echo "Use the restore.sh script to restore a commander state from a backup file."
    echo ""
    echo "Script requires following arguments:"
    echo "  1. commander directory path"
    echo "  2. badger database directory path"
    echo "  3. chain-spec directory path"
    echo "  4. geth chaindata directory path"
    echo "  5. path to the pigz tool"
    echo "  6. path to the pg_dump tool"
    echo "  7. aws bucket [optional]"
    echo ""
    echo "Example usage:"
    echo "bash $0 ./commander ./commander/db/badger/data ./commander/chain-spec ./commander/e2e/geth-data/geth/chaindata /usr/local/bin/pigz /usr/local/bin/pg_dump backup-aws-bucket"
    exit 0
fi

COMMANDER_DIR_PATH=$1
BADGER_DATA_DIR_PATH=$2
CHAIN_SPEC_DIR_PATH=$3
GETH_CHAINDATA_DIR_PATH=$4
PIGZ_PATH=$5
PG_DUMP_PATH=$6
AWS_S3_BUCKET=$7

# Prepare paths
BACKUP_DIR=$(date +"%Y-%m-%d_%H:%M:%S")
BACKUP_DIR_PATH="${COMMANDER_DIR_PATH}/backups/${BACKUP_DIR}"
COMPRESSED_BACKUP_DIR_PATH="${BACKUP_DIR_PATH}.tgz"

# Create a new backup directory based on the current time
mkdir -p "${BACKUP_DIR_PATH}"

# Backup badger data
rsync -a "${BADGER_DATA_DIR_PATH}"/* "${BACKUP_DIR_PATH}/badger/"

# Chain-spec data
rsync -a "${CHAIN_SPEC_DIR_PATH}"/* "${BACKUP_DIR_PATH}/chain-spec/"

# Backup geth chain data
rsync -a "${GETH_CHAINDATA_DIR_PATH}" "${BACKUP_DIR_PATH}/geth/"

# Dump postgres data
POSTGRES_IP=$(cut -d' ' -f1 <<<"$(hostname -I)") # Hardcode your machine IP here to be able to use this script on Mac
PGPASSWORD="password" "${PG_DUMP_PATH}" -h "${POSTGRES_IP}" -U root -p 5433 -C hubble -Fc -Z0 > "${BACKUP_DIR_PATH}/postgres.sql"

# Compress all the files
tar --use-compress-program="${PIGZ_PATH}" -cf "${COMPRESSED_BACKUP_DIR_PATH}" -C "${COMMANDER_DIR_PATH}" "./backups/${BACKUP_DIR}"

# Remove redundant uncompressed directory
rm -r "${BACKUP_DIR_PATH}"

# If AWS_S3_BUCKET is provided then upload the backup to a S3 bucket and remove the backup file from the local disk
if [[ -n ${AWS_S3_BUCKET:+x} ]]; then
    aws s3 cp "${COMPRESSED_BACKUP_DIR_PATH}" s3://"${AWS_S3_BUCKET}" >/dev/null
    rm -r "${COMPRESSED_BACKUP_DIR_PATH}"
fi
