# TodoList (Go + React + Postgres)

Features included:
- Devcontainer that runs as the host user (UID/GID passed via args)
- Go backend with CORS enabled and health endpoint
- React frontend built with Vite and served as static files by nginx
- Postgres DB and migrations (migrate/migrate image)
- docker-compose with healthchecks for DB, backend, and frontend

Quick steps:
docker ps |grep todolist-app
docker ps -a --filter "name=todolist-app" -q | xargs -r docker stop
docker ps -a --filter "name=todolist-app" -q | xargs -r docker rm

docker compose down
docker compose build
docker compose up -d
docker volume ls

### db
docker compose up db
docker compose down db -v
docker attach --detach-keys="ctrl-c" todolist-app-db-1
docker exec -it todolist-app-db-1 bash
pg_isready -U todo_user -d todo_db

### backend
docker compose build backend
docker compose build backend && docker compose run --rm backend sh
docker compose up backend -d
docker compose logs backend

docker compose down
docker compose down backend
docker compose down -v
http://localhost:8080/api/todos

### frontend
docker compose build frontend
docker compose build frontend && docker compose run --rm frontend sh
docker compose up frontend -d
docker compose logs frontend
docker compose down
docker compose down -v
http://localhost:3000


### set up .env file
```
echo "USER_UID=$(id -u)" >> .env
echo "USER_GID=$(id -g)" >> .env
DATABASE_URL=
POSTGRES_DB=
POSTGRES_PASSWORD=
POSTGRES_USER=
```