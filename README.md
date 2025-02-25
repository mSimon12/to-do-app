# To-Do App
A To-Do list Application to make management of tasks easier.


## Installation
To be able to run this application, some initial procedures need to be proceeded. First be sure that you have the following **requirements**:

- Go - [installing go](https://go.dev/doc/install)
- Docker - [installing docker](https://docs.docker.com/engine/install/)

With all the requirements ensured, we need to run a PostgreSQL container in docker, which will provide our application with a database server for storing the tasks list.

### Start database container
To run the database container, create a file **.env** at [database](api/src/models/database/) folder following the *example.env* to configure the PostgreSQL database. Than run the following command to start the DB docker container in a detached mode:
```
sudo docker compose up -d
```
To check that the container is running, call ``sudo docker ps`` and it should show something like this:
```
CONTAINER ID   IMAGE             COMMAND                  CREATED         STATUS          PORTS                                       NAMES
ab08b4b59c32   postgres:latest   "docker-entrypoint.s…"   3 minutes ago   Up 12 seconds   0.0.0.0:5432->5432/tcp, :::5432->5432/tcp   to_do_app_db
```

## Running
To run the **To-Do App** with only one command, we provide a bash script that will ensure all dependencies are installed, start the database container (if not running) and run the REST API server. For this, running the following command is enough:
```bash
bash app_init.sh
```