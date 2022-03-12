#!/bin/zsh

PRESIGN_URL='HERE'

# grab presigned url from latest backup and download into ./data/backup.tgz
curl -o ./data/backup.tgz $PRESIGN_URL
# unpack backup
tar --use-compress-program=unpigz -xf ./data/back.tgz -C ./data
