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

BACKUP_EXTENSION_LENGTH=$(echo "$1" | awk -F. '{print length($NF)}')
DECOMPRESSED_BACKUP_PATH=$(echo "$1" | rev | cut -c$((BACKUP_EXTENSION_LENGTH+2))- | rev)
POSTGRES_BACKUP_PATH=$DECOMPRESSED_BACKUP_PATH/postgres.sql
BADGER_BACKUP_PATH=$DECOMPRESSED_BACKUP_PATH/badger

# Decompress the compressed backup file
./unpigz.sh "$1"

# Restore postgres state
PGPASSWORD="$6" pg_restore -h "$3" -p "$4" -U "$5" -d "$7" -1 "$POSTGRES_BACKUP_PATH"

# Restore badger state
rsync -a "$BADGER_BACKUP_PATH" "$(dirname "$2")"

# Remove redundant decompressed backup directory
rm -rf "$DECOMPRESSED_BACKUP_PATH"
