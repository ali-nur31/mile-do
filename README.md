# Mile-Do

**Mile-Do** is a comprehensive task and goal management platform. This monorepo contains both the backend server and the frontend client (planned), orchestrated via Docker Compose.

## 📂 Repository Structure

```text
mile-do/
├── .env                         # Environment variables for Docker
├── client                       # Frontend Application (coming soon)
│   └── Dockerfile               # Dockerfile for frontend
├── docker-compose.yaml          # Main orchestration file for the entire stack
├── README.md
└── server                       # Backend API (Go, PostgreSQL, Redis)
    ├── cmd
    │   └── server
    │       └── main.go          # Main runnable file
    ├── config
    │   └── config.go            # Config vars loader
    ├── db
    │   ├── migrations           # Database migrations
    │   └── queries              # Sqlc queries to database
    ├── Dockerfile               # Dockerfile for server
    ├── docs                     # Swagger
    ├── internal
    │   ├── db                   # Generated database queries
    │   ├── domain               # Mapper for business logic
    │   ├── jobs                 # Background jobs
    │   ├── service              # Business logic
    │   └── transport
    │       └── http
    │           ├── middleware   # Middlewares for http transport endpoints
    │           └── v1           # Handlers
    │               └── dto      # Data transfer objects
    ├── pkg                      # External packages
    │   ├── asynq_jobs           # Client for background jobs
    │   ├── auth                 # Essentials for auth
    │   ├── logger               # Logger initialization
    │   ├── postgres             # Postgres client
    │   └── redis_db             # Redis client
    └── sqlc.yaml                # Sqlc config file
```

---

## 🚀 Quick Start

You can spin up the entire infrastructure (Database, Redis, Backend, Migrations) with a single command.

### Prerequisites

* [Docker](https://www.docker.com/products/docker-desktop/)
* [Docker Compose](https://docs.docker.com/compose/install/)

### Run the Project

1. **Configure Environment**
* Ensure you have a `.env` file in the root directory. You can copy the example from .env-example file in root directory

2. **Start Services**
```bash
docker-compose up -d --build
```

* `-d`: Detached mode (runs in background).
* `--build`: Forces a rebuild of images (useful after code changes).

3. **Check Status**
```bash
docker-compose ps
```

4. **View Logs**
```bash
docker-compose logs -f container_name
```

5. **Stop Services**
```bash
docker-compose down
```

---

## 🛠 Services & Ports

When running via Docker Compose, the following services are exposed:

| Service | Internal Host | External Port | Description |
| --- | --- |---------------| --- |
| **Server** | `server` | `8080`        | Main Go Backend API (`http://localhost:8080`) |
| **PostgreSQL** | `postgres` | `5435`        | Database (Use `localhost:5432` to connect via DBeaver) |
| **Redis** | `redis` | `6378`        | Task Queue & Caching |
| **Migrator** | `migrator` | N/A           | Ephemeral container. Runs on startup to apply SQL migrations. |

---

## 🐛 Troubleshooting Common Issues

**1. "port is already allocated"**

* Stop any local instances of Postgres or Redis running on your machine.
* Or change the external ports in `docker-compose.yml` (e.g., `"5435:5432"`).

**2. Database connection failed**

* Wait a few seconds. Postgres takes time to initialize on the first run. The `migrator` and `server` containers are configured to wait for it, but initialization might take longer on slower machines.

**3. Changes in code are not reflected**

* Docker does not watch for file changes by default. You must run `docker-compose up -d --build` to recompile the Go binary.