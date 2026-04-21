# To-Do App

![Docker](https://img.shields.io/badge/docker-ready-blue)
[![CI](https://github.com/mSimon12/to-do-app/actions/workflows/main.yml/badge.svg)](https://github.com/mSimon12/to-do-app/actions)
[![codecov](https://codecov.io/gh/mSimon12/to-do-app/branch/master/graph/badge.svg)](https://codecov.io/gh/mSimon12/to-do-app)

**To-Do App** is a full-stack task management application. It provides a RESTful backend API for creating, reading, updating, and deleting tasks, paired with a modern Angular web interface for a smooth user experience. Tasks are persisted in a PostgreSQL database and the entire system is containerized for easy deployment.

This application is composed of two services:
- A RESTful **Backend API** built with Go ([code here](./api))
- A modern **Angular Frontend** web interface ([code here](./web_interface))

## 🛠️ Tech Stack

<div align="left">
  <img src="https://skillicons.dev/icons?i=go" height="30" alt="go logo" />
  <img width="12" />
  <img src="https://skillicons.dev/icons?i=typescript" height="30" alt="typescript logo" />
  <img width="12" />
  <img src="https://skillicons.dev/icons?i=angular" height="30" alt="angular logo" />
  <img width="12" />
  <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/postgresql/postgresql-original.svg" height="30" alt="postgresql logo" />
  <img width="12" />
  <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/docker/docker-original.svg" height="30" alt="docker logo" />
</div>

---

## 🚀 Features

- 📋 Full CRUD operations for task management
- 🔍 Query single or multiple tasks
- 🎨 Modern, responsive Angular web interface
- 🐳 Fully Dockerized with multi-service support
- 📖 Auto-generated API documentation endpoint
- 🧪 Unit testing and basic CI setup

---

## ⚙️ Getting Started

### Requirements

All you need is Docker installed on your machine, since the entire application is containerized.

- Docker — [installing Docker](https://docs.docker.com/engine/install/)

### 🔧 Clone the Repository

```bash
git clone https://github.com/mSimon12/to-do-app
cd to-do-app
```

### 🔐 Environment Setup

Create a **.env** file inside the [deploy](deploy) folder by copying the provided example and filling in your PostgreSQL credentials:

```bash
cp deploy/example.env deploy/.env
# Edit deploy/.env with your database credentials
```

---

## ▶️ Running the Application

You can run the system either directly on your machine or using Docker.

### 🐳 Option 1: Run with Docker (Recommended)

Start all services (database, API, and web interface) using Docker Compose:

```bash
docker compose -f deploy/docker-compose.yml up -d --build
```

> It might take a few minutes to initialize the PostgreSQL database for the first time.

Once running, access the services at:
- **Web Interface**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **API Documentation**: http://localhost:8080/api/documentation/

To verify all containers are up:

```bash
docker ps
```

Expected output:

```
CONTAINER ID   IMAGE                COMMAND               CREATED        STATUS        PORTS                    NAMES
xxxxxxxxxxxx   to_do_app_ui:latest  "..."                 1 min ago      Up 1 min      0.0.0.0:3000->4000/tcp   to_do_app_ui
xxxxxxxxxxxx   to_do_app_api:latest "/docker-to-do-api"  1 min ago      Up 1 min      0.0.0.0:8080->8080/tcp   to_do_app_api
xxxxxxxxxxxx   postgres:latest      "docker-entrypoint…" 1 min ago      Up 1 min      0.0.0.0:5432->5432/tcp   to_do_app_db
```

### 💻 Option 2: Run Without Containers

Open two terminals.

**Terminal 1 — Backend API**

```bash
cd api
go run .
```

**Terminal 2 — Angular Frontend**

```bash
cd web_interface
npm install
npm start
```

The web interface will be available at `http://localhost:4200`.

---

## 🛣️ Roadmap

### 🔧 Backend

- [x] Built a REST API using Go
- [x] Configured PostgreSQL database integration
- [x] Implemented CRUD operations for tasks
- [x] Auto-generated API documentation endpoint
- [x] Dockerized backend service

### 🎨 Frontend

- [x] Built a modern Single Page Application (SPA) using Angular with TypeScript
- [x] Implemented task listing board (`main-board` component)
- [x] Created task card components for individual task display
- [x] Built task detail view with full task information
- [x] Connected frontend to backend via HTTP service layer
- [x] Added theme service for UI customization
- [x] Dockerized frontend service

---

## 🧩 System Architecture

The To-Do App is composed of two isolated services — a **Go Backend API** and an **Angular Frontend** — both deployed as independent containers alongside a **PostgreSQL** database.

### 🔩 Backend (API Service)

The backend is a Go REST API that exposes endpoints for managing tasks. It connects directly to a PostgreSQL database and handles all business logic and data persistence.

### 🎨 Frontend (Web Interface)

The frontend is a modern **Single Page Application (SPA)** built with **Angular 21** and **TypeScript**. It includes:

- **Components**: Reusable UI building blocks (`main-board`, `task-card`, `task-details`) with encapsulated styles and logic.
- **Services**: HTTP abstraction layer (`tasks-api`) that handles communication with the backend API, and a `theme` service for UI theming.
- **Routes**: Angular routing enables client-side navigation between views without full page reloads.

### 🔁 Request Flow

1. The user interacts with the **Angular Frontend**.
2. **Angular Services** prepare and send HTTP requests to the Go backend API.
3. The **Backend API** processes the request and reads/writes to the **PostgreSQL** database.
4. The **API Response** is returned to the frontend service, which updates component state and triggers UI re-rendering.
5. The **Angular component** displays the updated data to the user in real-time.
