#!/bin/bash

## PostgreSQL
echo Starting To-Do App Database
if sudo lsof -i :5432; then
    echo Stopping local postgres
    sudo service postgresql stop
fi
sudo docker compose -f api/deploy/docker-compose.yml up -d

## Application
echo Running To-Do App
## Run api in Go
go run .
