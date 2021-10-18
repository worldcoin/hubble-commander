#!/bin/bash
if [ "$#" -ne 1 ]; then
    echo "Script used for decompressing the backup .tgz file."
    echo "The decompressed backup will be put in the same directory as .tgz file."
    echo ""
    echo "Script requires 1 argument to work:"
    echo "  1. backup .tgz file path"
    echo ""
    echo "Example usage:"
    echo "bash $0 ./commander/backups/2021-10-05_13:20:38.tgz"
    exit 0
fi

COMPRESSED_BACKUP_DIR_PATH=$1

DECOMPRESSED_TARGET_DIR=$(dirname "${COMPRESSED_BACKUP_DIR_PATH}")
tar --use-compress-program=unpigz -xf "${COMPRESSED_BACKUP_DIR_PATH}" -C "${DECOMPRESSED_TARGET_DIR}"
