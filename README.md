# To-Do App

To-Do App is a REST API developed in Go that provides a system for managing tasks to be done. It provides the basic functions of Creating, Updating, Deleting and Querying single or multiple tasks.
The tasks are stored in a PostgreSQL database and it have been developed to be put in production with very little commands required.


## Requirements
To be able to run this application, all you will need is to have docker installed, since the whole application is containerized.

- Docker - [installing docker](https://docs.docker.com/engine/install/)

## Running

### Configure database container
To configure the database container, create a file **.env** at [database](api/deploy) folder following the *example.env* and update the configuration for your PostgreSQL database. 

### Create and Run containers

This application was developed to be run as easy as possible, and this can be accomplished by simply using the docker compose to build the required containers.

```bash
sudo docker compose -f docker-compose.yml up -d --build
```

The command above should be enough to run the system. At this point you should have 2 containers running in the background. One container responsible for the PostgreSQL database and another running the application and providing the API endpoints.
To be sure that both are running, you can check if the containers were correctly created by calling  ``sudo docker ps`` and it should show something like this:

```                                 NAMES
CONTAINER ID   IMAGE             COMMAND                  CREATED         STATUS         PORTS                                         NAMES
22260beeadc6   to_do_app_image   "/docker-to-do-api"      3 minutes ago   Up 3 minutes   0.0.0.0:8080->8080/tcp, [::]:8080->8080/tcp   to_do_app
f9ebee888816   postgres:latest   "docker-entrypoint.s…"   3 minutes ago   Up 3 minutes   0.0.0.0:5432->5432/tcp, [::]:5432->5432/tcp   to_do_app_db
```
