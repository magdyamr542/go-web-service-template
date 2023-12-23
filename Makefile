include .env
export $(shell sed 's/=.*//' .env)

GIT_SHA := $(shell git rev-parse HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT_TIME := $(shell git show --no-patch --format=%ci HEAD | tr -s ' ' '|')
BUILD_TIME := $(shell date +'%Y-%m-%d %H:%M:%S %z' | tr -s ' ' '|')
VERSION := "1.0.0"
CWD := $(shell pwd)
NEW_MIGRATION_COMMAND := sql-migrate new -env=development
APPLY_MIGRATIONS_COMMAND := sql-migrate up -env=development

.PHONY: build
build:
	@./scripts/build.sh

.PHONY: generate
generate:
	oapi-codegen \
	-generate types,server,spec \
	-package api openapi/api.yaml > pkg/api/api.gen.go

.PHONY: run
run: build
	@./binary -environment development --with-metrics

.PHONY: new-migration
new-migration:
	@cd ./pkg/storage/migrations && $(NEW_MIGRATION_COMMAND)

.PHONY: migrate-up
migrate-up:
	@cd ./pkg/storage/migrations && $(APPLY_MIGRATIONS_COMMAND)

.PHONY: connect-db
connect-db:
	docker compose exec -it db psql -W app -U user


.PHONY: docker-up
docker-up:
	@docker compose --profile monitoring up

.PHONY: seed-db
seed-db:
	@pkg/storage/seeder.sh