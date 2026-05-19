#!/bin/bash

set -e

MOCKGEN="$(go env GOPATH)/bin/mockgen"

if [ ! -f "$MOCKGEN" ]; then
    echo "Installing mockgen..."
    go install go.uber.org/mock/mockgen@v0.6.0
fi

generate_mock() {
    local service=$1
    local type=$2
    local src=$3
    local dest=$4

    echo "Generating mock for $service $type..."
    "$MOCKGEN" -source="$src" -destination="$dest"
    echo "Completed mock generation for $service $type"
}

services=("syllabus" "user-management")

for service in "${services[@]}"; do
    generate_mock "$service" "service"    "internal/$service/service/interface.go"    "internal/$service/service/mocks/mock.go" &
    generate_mock "$service" "repository" "internal/$service/repository/interface.go" "internal/$service/repository/mocks/mock.go" &
done

wait

echo "All mock generations completed."
