#!/bin/bash
if [ "$#" -ne 1 ]; then
    echo "Script used for decompressing the backup .tgz file."
    echo "The decompressed backup will be put in the same directory as .tgz file."
    echo ""
    echo "Script requires 1 argument to work:"
    echo "  1. backup .tgz file path"
    echo ""
    echo "Example usage:"
    echo "$0 ./commander/backups/2021-10-05_13:20:38.tgz"
    exit 0
fi

DECOMPRESSED_TARGET_DIR=$(dirname "$(dirname "$1")")
tar --use-compress-program=unpigz -xf "$1" -C "$DECOMPRESSED_TARGET_DIR"
