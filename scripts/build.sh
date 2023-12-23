#!/bin/sh

GIT_SHA=$(git rev-parse HEAD)
GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
GIT_COMMIT_TIME=$(git show --no-patch --format=%ci HEAD | tr -s ' ' '|')
BUILD_TIME=$(date +'%Y-%m-%d %H:%M:%S %z' | tr -s ' ' '|')
VERSION="1.0.0"

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/magdyamr542/go-web-service-template/pkg/version.Version=$VERSION -X github.com/magdyamr542/go-web-service-template/pkg/version.BuildTime=$BUILD_TIME -X github.com/magdyamr542/go-web-service-template/pkg/version.CommitTime=$GIT_COMMIT_TIME -X github.com/magdyamr542/go-web-service-template/pkg/version.CommitSHA=$GIT_SHA -X github.com/magdyamr542/go-web-service-template/pkg/version.CommitBranch=$GIT_BRANCH" \
	-o binary
