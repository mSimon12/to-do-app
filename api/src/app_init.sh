#!/bin/bash


## Install dependencies
go get -u github.com/gin-gonic/gin
go get github.com/jackc/pgx/v5
go get github.com/joho/godotenv

## PostgreSQL

echo Starting To-Do App Database
sudo docker compose -f models/database/docker-compose.yml up -d


echo Running To-Do App
## Run api in Go
go run .
