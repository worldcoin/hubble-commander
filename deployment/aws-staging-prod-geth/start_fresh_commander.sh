#!/bin/bash
docker-compose down -v
docker-compose up mining-geth-init
docker-compose up -d mining-geth
docker-compose up api-geth-init
docker-compose up -d api-geth
docker-compose up deploy
docker-compose up -d commander
