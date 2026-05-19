#!/bin/bash

CONTAINER_NAME=$1

echo "Waiting for Postgres container '$CONTAINER_NAME' to be ready..."

max_attempts=10
attempt=0

while [ $attempt -lt $max_attempts ]; do
    if docker exec $CONTAINER_NAME pg_isready -U postgres -d products > /dev/null 2>&1; then
        echo "Postgres container '$CONTAINER_NAME' is ready!"
        exit 0
    fi

    attempt=$((attempt + 1))
    echo "Attempt $attempt/$max_attempts: Postgres not ready yet..."
    sleep 2
done

echo "ERROR: Postgres container '$CONTAINER_NAME' did not become ready in time"
exit 1