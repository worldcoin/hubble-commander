#!/bin/bash

if [ "$#" -ne 9 ]; then
    echo "Script used for restoring a commander state from a backup file."
    echo "This script requires pg_restore tool to work."
    echo ""
    echo "Script requires 9 arguments:"
    echo "  1. unpigz.sh script path"
    echo "  2. backups .tgz file path"
    echo "  3. badger database directory path"
    echo "  4. chain-spec.yaml file path"
    echo "  5. postgres host"
    echo "  6. postgres port"
    echo "  7. postgres user"
    echo "  8. postgres password"
    echo "  9. postgres dbname"
    echo ""
    echo "Example usage:"
    echo "bash $0 ./deployment/scripts/unpigz.sh ./backups/2021-10-05_16:00:19.tgz ./deployment/db/badger ./deployment/chain-spec.yaml 192.168.0.81 5433 root password hubble"
    exit 0
fi

UNPIGZ_SCRIPT_PATH=$1
COMPRESSED_BACKUP_DIR_PATH=$2
BADGER_DATA_DIR_PATH=$3
CHAIN_SPEC_FILE_PATH=$4
POSTGRES_IP=$5
POSTGRES_PORT=$6
POSTGRES_USER=$7
POSTGRES_PASSWORD=$8
POSTGRES_DBNAME=$9

# These 2 lines cut the extension from COMPRESSED_BACKUP_DIR_PATH
BACKUP_EXTENSION_LENGTH=$(echo "${COMPRESSED_BACKUP_DIR_PATH}" | awk -F. '{print length($NF)}')
DECOMPRESSED_BACKUP_PATH=$(echo "${COMPRESSED_BACKUP_DIR_PATH}" | rev | cut -c$((BACKUP_EXTENSION_LENGTH+2))- | rev)

# Prepare paths
POSTGRES_BACKUP_PATH=$DECOMPRESSED_BACKUP_PATH/postgres.sql
BADGER_BACKUP_PATH=$DECOMPRESSED_BACKUP_PATH/badger/
CHAIN_SPEC_BACKUP_PATH=$DECOMPRESSED_BACKUP_PATH/chain-spec.yaml

# Decompress the compressed backup file using the unpigz.sh script
bash "${UNPIGZ_SCRIPT_PATH}" "${COMPRESSED_BACKUP_DIR_PATH}"

# Restore postgres state
PGPASSWORD="${POSTGRES_PASSWORD}" pg_restore -h "${POSTGRES_IP}" -p "${POSTGRES_PORT}" -U "${POSTGRES_USER}" -d "${POSTGRES_DBNAME}" -1 "${POSTGRES_BACKUP_PATH}"

# Restore badger state
rsync -a "${BADGER_BACKUP_PATH}" "${BADGER_DATA_DIR_PATH}"

# Restore chain-spec data
rsync -a "${CHAIN_SPEC_BACKUP_PATH}" "${CHAIN_SPEC_FILE_PATH}"

# Remove redundant decompressed backup directory
rm -r "${DECOMPRESSED_BACKUP_PATH}"
