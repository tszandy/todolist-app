# Project structure

This document describes the repository layout, main responsibilities of folders/files, and common commands to build and run the project. Updated: December 24, 2025.

## Overview

A todo-list web application composed of a Go backend, a React frontend (JSX, Vite), Docker artifacts, and SQL migrations. The repository uses Docker Compose to run services together.

## Top-level layout

```
/ (project root)
├── docker-compose.yml            # Compose file to run services (backend, frontend, db, migrations)
├── docker-compose.dev.yml        # Dev-only compose file (hot reload, bind mounts)
├── docker-compose.prod.yml       # Production compose file (optimized images)
├── README.md
├── structure.md                  # This file
├── backend/                      # Go backend service
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   └── tmp/                      # (build output, optional)
│       └── main
├── frontend/                     # React frontend (Vite)
│   ├── Dockerfile
│   ├── index.html
│   ├── nginx.conf
│   ├── package.json
│   ├── package-lock.json
│   └── src/
│       ├── App.jsx
│       └── main.jsx
└── migrations/                   # SQL migration files
    ├── 1_init.up.sql
    └── 1_init.down.sql
```

## Folder / file descriptions

- `docker-compose.yml`, `docker-compose.dev.yml`, `docker-compose.prod.yml`
  - Orchestrate services for local development and production. Services: backend, frontend, db, migrations. Dev/prod files use different build targets and volumes.

- `backend/`
  - Go-based HTTP API service. Entrypoint: `main.go`. Uses environment variables: `DATABASE_URL`, `PORT` (default 8080).
  - `Dockerfile` — multi-stage build for dev/prod. Dev uses bind mounts and hot reload; prod is optimized.
  - `tmp/main` — build output (optional, ignored in most workflows).

- `frontend/`
  - React app (JSX, Vite). Entrypoint: `src/App.jsx`, `src/main.jsx`. Uses Vite scripts: `dev`, `build`, `preview`.
  - `nginx.conf` — used to serve static assets in production.
  - `Dockerfile` — builds and serves frontend for Docker Compose.

- `migrations/`
  - SQL migrations for database schema. `1_init.up.sql` creates the `todos` table; `1_init.down.sql` drops it.

- `README.md`
  - Project README. Check for setup, environment variables, and usage notes.

## Common workflows

- **Start everything via Docker Compose (recommended):**

  ```zsh
  docker-compose up --build
  docker-compose down
  ```

- **Backend local development**

  ```zsh
  cd backend
  go mod download
  go run .
  ```
  - For hot reload: use `docker-compose.dev.yml` (bind mounts, dev target).

- **Frontend local development**

  ```zsh
  cd frontend
  npm install
  npm run dev
  npm run build   # for production bundle
  npm run preview # preview production build
  ```

- **Database migrations**

  - Migrations are in `migrations/`. Run with the `migrate` service in Docker Compose, or apply SQL files directly.

## Configuration and environment

- Environment variables: see `.env.example` (recommended to create), `docker-compose.yml`, and `backend/main.go` for required keys (`DATABASE_URL`, `PORT`, etc).
- Database defaults: user `todo_user`, password `todo_pass`, db `todo_db` (see `docker-compose.yml`).

## API endpoints (backend)

- `GET /api/todos` — list todos
- `POST /api/todos` — create todo
- `PUT /api/todos/{id}` — toggle completed
- `GET /api/health` — health check

## Useful commands (macOS, zsh)

- Build & run with Docker Compose:
  ```zsh
  docker-compose up --build
  docker-compose down
  ```
- Backend local run:
  ```zsh
  cd backend && go mod download && go run .
  ```
- Frontend local run:
  ```zsh
  cd frontend && npm install && npm run dev
  ```

## Notes and recommendations

- Add a `.env.example` in the project root listing required environment variables (DB host, port, username, password, JWT secrets, etc.).
- Document any additional services (database type/version) that `docker-compose.yml` expects.
- Keep `migrations/` in version control and apply them consistently in CI/CD or container startup.

---

For exact run commands, inspect `docker-compose.yml`, `backend/main.go`, and `frontend/package.json` for up-to-date values.
