#!/bin/zsh

# checkeout branch with modified export
git checkout migrate-0.4.0
# build cli
go build -o build/hubble-cli ./main

./build/hubble-cli export -type state --file state_dump.json
./build/hubble-cli export -type accounts --file accounts_dump.json
