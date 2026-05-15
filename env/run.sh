export GIT_COMMIT=$(git rev-parse HEAD)
echo "Запуск Docker Compose с GIT_COMMIT=${GIT_COMMIT}..."

docker compose up
