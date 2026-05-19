-include .env
export

PREFIX  = edu-planner-
NETWORK = edu-planner_default

DOCKER_COMPOSE  = docker compose
DIR             = $(CURDIR)
GOLINT_VERSION  = v2.12.2

# Map service name to its postgres database name
DB_NAME = $(if $(filter user-management,$(SERVICE)),user_management,$(SERVICE))
DB_URL  = postgres://$(DB_USER):$(DB_PASSWORD)@edu-planner-$(SERVICE)-postgres:5432/$(DB_NAME)?sslmode=disable

.PHONY: up down build service-start service-rebuild migrations-up migrations-down go-unit-tests go-integration-tests go-test generate-mocks db-dump-up

up:
	$(DOCKER_COMPOSE) up -d edu-planner-user-management-postgres edu-planner-syllabus-postgres localstack
	./hack/scripts/run-migrations.sh
	$(DOCKER_COMPOSE) up -d

down:
	$(DOCKER_COMPOSE) down

build:
	$(DOCKER_COMPOSE) build

service-start:
	$(DOCKER_COMPOSE) up -d edu-planner-$(SERVICE)

service-rebuild:
	$(DOCKER_COMPOSE) build edu-planner-$(SERVICE)
	$(DOCKER_COMPOSE) up -d --no-deps edu-planner-$(SERVICE)

.PHONY:migrations-up
migrations-up:
	docker run --rm -v $(CURDIR)/internal/$(patsubst ${PREFIX}%,%,${SERVICE})/repository/postgres/migration:/migration \
		--network $(NETWORK) \
		migrate/migrate:v4.15.2 -path=/migration \
		-database "postgres://${DB_USER}:${DB_PASSWORD}@${PREFIX}${SERVICE}-postgres:5432/$(DB_NAME)?sslmode=disable" up

.PHONY:migrations-down
migrations-down:
	docker run --rm -v $(CURDIR)/internal/$(patsubst ${PREFIX}%,%,${SERVICE})/repository/postgres/migration:/migration \
		--network $(NETWORK) \
		migrate/migrate:v4.15.2 -path=/migration \
		-database "postgres://${DB_USER}:${DB_PASSWORD}@${PREFIX}${SERVICE}-postgres:5432/$(DB_NAME)?sslmode=disable" down -all

go-unit-tests:
	./hack/scripts/run-go-unit-tests.sh

go-integration-tests:
	./hack/scripts/run-go-integration-tests.sh

go-test:
	./hack/scripts/run-go-unit-tests.sh
	./hack/scripts/run-go-integration-tests.sh

generate-mocks:
	./hack/scripts/generate-mocks.sh

# Load demo seed data into both databases.
# Run AFTER: make up  (containers must be healthy and migrations applied)
db-dump-up:
	docker exec -i edu-planner-user-management-postgres psql -U $(DB_USER) -d user_management < hack/dump/user-management.sql
	docker exec -i edu-planner-syllabus-postgres      psql -U $(DB_USER) -d syllabus       < hack/dump/syllabus.sql

.PHONY:go-lint
go-lint:
	docker run -t --rm -v $(DIR):/app -v ~/.cache/golangci-lint/$(GOLINT_VERSION):/root/.cache -w /app golangci/golangci-lint:$(GOLINT_VERSION) golangci-lint run -v
