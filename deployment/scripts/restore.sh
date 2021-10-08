#!/bin/bash

if [ "$#" -ne 7 ]; then
    echo "Script used for restoring a commander state from a backup file."
    echo ""
    echo "Script requires 7 arguments:"
    echo "  1. backups .tgz file path"
    echo "  2. badger database directory path"
    echo "  3. postgres host"
    echo "  4. postgres port"
    echo "  5. postgres user"
    echo "  6. postgres password"
    echo "  7. postgres dbname"
    echo ""
    echo "Example usage:"
    echo "$0 ./backups/2021-10-05_16:00:19.tgz ./deployment/db/badger 192.168.0.81 5433 root password hubble"
    exit 0
fi

COMPRESSED_BACKUP_DIR_PATH=$1
BADGER_DATA_DIR_PATH=$2
POSTGRES_IP=$3
POSTGRES_PORT=$4
POSTGRES_USER=$5
POSTGRES_PASSWORD=$6
POSTGRES_DBNAME=$7

# These 2 lines cut the extension from COMPRESSED_BACKUP_DIR_PATH
BACKUP_EXTENSION_LENGTH=$(echo "${COMPRESSED_BACKUP_DIR_PATH}" | awk -F. '{print length($NF)}')
DECOMPRESSED_BACKUP_PATH=$(echo "${COMPRESSED_BACKUP_DIR_PATH}" | rev | cut -c$((BACKUP_EXTENSION_LENGTH+2))- | rev)

# Prepare paths
POSTGRES_BACKUP_PATH=$DECOMPRESSED_BACKUP_PATH/postgres.sql
BADGER_BACKUP_PATH=$DECOMPRESSED_BACKUP_PATH/badger

# Decompress the compressed backup file
./unpigz.sh "${COMPRESSED_BACKUP_DIR_PATH}"

# Restore postgres state
PGPASSWORD="${POSTGRES_PASSWORD}" pg_restore -h "${POSTGRES_IP}" -p "${POSTGRES_PORT}" -U "${POSTGRES_USER}" -d "${POSTGRES_DBNAME}" -1 "$POSTGRES_BACKUP_PATH"

# Restore badger state
rsync -a "$BADGER_BACKUP_PATH" "$(dirname "${BADGER_DATA_DIR_PATH}")"

# Remove redundant decompressed backup directory
rm -rf "$DECOMPRESSED_BACKUP_PATH"
