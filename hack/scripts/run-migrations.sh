#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Load .env
if [ -f "$ROOT_DIR/.env" ]; then
    set -a
    # shellcheck source=/dev/null
    source "$ROOT_DIR/.env"
    set +a
fi

DB_USER="${DB_USER:-dev}"
DB_PASSWORD="${DB_PASSWORD:-12345}"
NETWORK="${NETWORK:-edu-planner_default}"
MIGRATE_IMAGE="migrate/migrate:v4.15.2"

wait_healthy() {
    local container="$1"
    echo "Waiting for $container to be healthy..."
    until [ "$(docker inspect "$container" --format='{{.State.Health.Status}}' 2>/dev/null)" = "healthy" ]; do
        sleep 2
    done
    echo "$container is healthy."
}

wait_healthy edu-planner-user-management-postgres
wait_healthy edu-planner-syllabus-postgres

echo "Running user-management migrations..."
docker run --rm \
    -v "$ROOT_DIR/internal/user-management/repository/postgres/migration:/migration" \
    --network "$NETWORK" \
    "$MIGRATE_IMAGE" \
    -path=/migration \
    -database "postgres://$DB_USER:$DB_PASSWORD@edu-planner-user-management-postgres:5432/user_management?sslmode=disable" \
    up

echo "Running syllabus migrations..."
docker run --rm \
    -v "$ROOT_DIR/internal/syllabus/repository/postgres/migration:/migration" \
    --network "$NETWORK" \
    "$MIGRATE_IMAGE" \
    -path=/migration \
    -database "postgres://$DB_USER:$DB_PASSWORD@edu-planner-syllabus-postgres:5432/syllabus?sslmode=disable" \
    up

echo "All migrations complete."
