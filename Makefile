GIT_SHA := $(shell git rev-parse HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT_TIME := $(shell git show --no-patch --format=%ci HEAD | tr -s ' ' '|')
BUILD_TIME := $(shell date +'%Y-%m-%d %H:%M:%S %z' | tr -s ' ' '|')
VERSION := "1.0.0"

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