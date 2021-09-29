#!/bin/bash
docker-compose down -v
docker-compose up deploy
docker-compose up -d commander
