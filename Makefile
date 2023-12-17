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
	@go build -ldflags "-X github.com/magdyamr542/go-web-service-template/pkg/handler.Version=$(VERSION) -X github.com/magdyamr542/go-web-service-template/pkg/handler.BuildTime=$(BUILD_TIME) -X github.com/magdyamr542/go-web-service-template/pkg/handler.CommitTime=$(GIT_COMMIT_TIME) -X github.com/magdyamr542/go-web-service-template/pkg/handler.CommitSHA=$(GIT_SHA) -X github.com/magdyamr542/go-web-service-template/pkg/handler.CommitBranch=$(GIT_BRANCH)" \
		-o binary

.PHONY: generate
generate:
	oapi-codegen \
	-generate types,server,spec \
	-package api openapi/api.yaml > pkg/api/api.gen.go

.PHONY: run
run: build
	@./binary -environment development 

.PHONY: sqlc-generate
sqlc-generate:
	@docker run --rm  \
		--network go-web-service-template_default \
		-v $(CWD)/pkg/storage:/src \
		-w /src/sqlc \
		-e DATABASE_HOST=db \
		-e DATABASE_PORT=$(DATABASE_PORT) \
		-e DATABASE_NAME=$(DATABASE_NAME) \
		-e DATABASE_USER=$(DATABASE_USER) \
		-e DATABASE_PASSWORD=$(DATABASE_PASSWORD) \
		sqlc/sqlc:1.24.0 generate

.PHONY: new-migration
new-migration:
	@cd ./pkg/storage/migrations && $(NEW_MIGRATION_COMMAND)

.PHONY: migrate-up
migrate-up:
	@cd ./pkg/storage/migrations && $(APPLY_MIGRATIONS_COMMAND)

.PHONY: connect-db
connect-db:
	docker compose exec -it db psql -W app -U user
