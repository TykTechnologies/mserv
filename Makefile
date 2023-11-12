SHELL := bash

# Default - top level rule is what gets run when you run just `make` without specifying a goal/target.
.DEFAULT_GOAL := build

.DELETE_ON_ERROR:
.ONESHELL:
.SHELLFLAGS := -euo pipefail -c

MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --warn-undefined-variables


ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later.)
endif
.RECIPEPREFIX = >

binary_name ?= $(shell basename $(CURDIR))
image_repository ?= tykio/$(binary_name)

# Adjust the width of the first column by changing the '16' value in the printf pattern.
help:
> @grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
  | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'
.PHONY: help

all: test lint build build-cli ## Test and lint and build.
build: out/image-id ## Build the mserv Docker image. Will also test and lint.
build-binary: $(binary_name) ## Build a bare binary only, without a Docker image wrapped around it.
build-cli: mservctl/mservctl ## Build the mservctl CLI binary. Will also test and lint.
.PHONY: all test lint build build-binary build-cli

check-swagger:
> which swagger || (GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger)

swagger: check-swagger
> GO111MODULE=on go mod vendor
> GO111MODULE=off swagger generate spec -o ./doc/swagger.yaml --scan-models -x mservclient -x vendor
.PHONY: swagger

serve-swagger: check-swagger
> swagger serve -F=swagger ./doc/swagger.yaml

swagger-client: check-swagger
> mkdir -p ./mservclient
> swagger generate client -f ./doc/swagger.yaml -t ./mservclient

clean: ## Clean up the temp and output directories, and any built binaries. This will cause everything to get rebuilt.
> rm -rf ./tmp ./out
> go clean
> cd mservctl
> go clean
.PHONY: clean

clean-docker: ## Clean up any built Docker images.
> docker images \
  --filter=reference=$(image_repository) \
  --no-trunc --quiet | sort -f | uniq | xargs -n 1 docker rmi --force
> rm -f out/image-id
.PHONY: clean-docker

clean-hack: ## Clean up binaries under 'hack'.
> rm -rf ./hack/bin
.PHONY: clean-hack

clean-all: clean clean-docker clean-hack ## Clean all of the things.
.PHONY: clean-all

# Run go tests
test: $(shell find . -type f -iname "*.go")
> mkdir -p $(@D)
> go test -v -count=1 -p 1 -race ./...

# Lint golangci lint
lint: .golangci.yaml hack/bin/golangci-lint
> mkdir -p $(@D)
> hack/bin/golangci-lint run

hack/bin/golangci-lint:
> curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
> | sh -s -- -b $(shell pwd)/hack/bin

# Docker image - re-build if the lint output is re-run.
out/image-id: Dockerfile
> mkdir -p $(@D)
> image_id="$(image_repository):$(shell uuidgen)"
> DOCKER_BUILDKIT=1 docker build --tag="$${image_id}" .
> echo "$${image_id}" > out/image-id

# Main server binary
$(binary_name):
> go build -mod=vendor

# CLI binary - re-build if the lint output is re-run.
mservctl/mservctl:
> cd mservctl
> go build -mod=vendor

# Start runs development environment with mserv and mongo in docker-compose
start:
> docker-compose up -d

# Stop runs development environment with mserv and mongo in docker-compose
stop:
> docker-compose stop
