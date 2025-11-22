# TodoList (Go + React + Postgres)

Features included:
- Devcontainer that runs as the host user (UID/GID passed via args)
- Go backend with CORS enabled and health endpoint
- React frontend built with Vite and served as static files by nginx
- Postgres DB and migrations (migrate/migrate image)
- docker-compose with healthchecks for DB, backend, and frontend

Quick steps:
1. Copy `.env.example` to `.env` and edit if needed.
2. Start services: `docker compose up --build`
   - The `migrate` service will run once at startup to apply migrations.
3. Visit http://localhost:3000

