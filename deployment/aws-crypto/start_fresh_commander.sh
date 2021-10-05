#!/bin/bash
docker-compose down -v
sudo rm -rf ./badger-data
sudo rm -rf ./geth-data/geth
docker-compose up deploy
docker-compose up -d commander
