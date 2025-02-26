#!/bin/bash


## Install dependencies
go get -u github.com/gin-gonic/gin
go get github.com/jackc/pgx/v5
go get github.com/joho/godotenv

## PostgreSQL
echo Starting To-Do App Database
if sudo lsof -i :5432; then
    echo Stopping local postgres
    sudo service postgresql stop
fi
sudo docker compose -f models/database/docker-compose.yml up -d

## Application
echo Running To-Do App
## Run api in Go
go run .
