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
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
docker volume ls

### set up .env file
```
echo "USER_UID=$(id -u)" >> .env
echo "USER_GID=$(id -g)" >> .env
echo "POSTGRES_USER=todo_user" >> .env
echo "POSTGRES_PASSWORD=todo_pass" >> .env
echo "POSTGRES_DB=todo_db" >> .env
echo "DATABASE_URL=postgres://todo_user:todo_pass@db:5432/todo_db?sslmode=disable" >> .env
```

### all prod
docker compose -f docker-compose.yml -f docker-compose.prod.yml up
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
docker compose -f docker-compose.yml -f docker-compose.prod.yml down
docker compose -f docker-compose.yml -f docker-compose.prod.yml down -v
### all dev
docker compose -f docker-compose.yml -f docker-compose.dev.yml up
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d
docker compose -f docker-compose.yml -f docker-compose.dev.yml down
docker compose -f docker-compose.yml -f docker-compose.dev.yml down -v


### db
docker compose -f docker-compose.yml -f docker-compose.dev.yml up db
docker compose -f docker-compose.yml -f docker-compose.dev.yml down db
docker compose -f docker-compose.yml -f docker-compose.dev.yml down db -v
docker attach --detach-keys="ctrl-c" todolist-app-db-1
docker exec -it todolist-app-db-1 bash
pg_isready -U todo_user -d todo_db

### backend
docker compose -f docker-compose.yml -f docker-compose.dev.yml build backend
docker compose -f docker-compose.yml -f docker-compose.dev.yml up backend -d
docker compose -f docker-compose.yml -f docker-compose.dev.yml run backend sh
docker compose -f docker-compose.yml -f docker-compose.dev.yml logs backend
docker compose -f docker-compose.yml -f docker-compose.dev.yml down backend
docker compose -f docker-compose.yml -f docker-compose.dev.yml down backend -v
http://localhost:8080/api/todos

### frontend
docker compose -f docker-compose.yml -f docker-compose.dev.yml build frontend
docker compose -f docker-compose.yml -f docker-compose.dev.yml up frontend -d
docker compose -f docker-compose.yml -f docker-compose.dev.yml logs frontend
docker compose -f docker-compose.yml -f docker-compose.dev.yml down frontend
docker compose -f docker-compose.yml -f docker-compose.dev.yml down frontend -v
http://localhost:3000

