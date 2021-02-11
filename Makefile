SHELL := bash

# Default - top level rule is what gets run when you run just `make` without specifying a goal/target.
.DEFAULT_GOAL := build

.DELETE_ON_ERROR:
.ONESHELL:
.SHELLFLAGS := -euo pipefail -c

MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --warn-undefined-variables

IMAGE_NAME ?= tykio/mserv

ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later.)
endif
.RECIPEPREFIX = >

# Adjust the width of the first column by changing the '16' value in the printf pattern.
help:
> @grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
  | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'
.PHONY: help

all: test lint build ## Test and lint and build.
test: tmp/.tests-passed.sentinel ## Run tests.
lint: tmp/.linted.sentinel ## Lint all of the Go code. Will also test.
build: out/image-id ## Build the mserv Docker image. Will also test and lint.
build-cli: mservctl/mservctl ## Build the mservctl CLI binary. Will also test and lint.
.PHONY: all test lint build build-cli

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
  --filter=reference=$(IMAGE_NAME) \
  --no-trunc --quiet | sort -f | uniq | xargs -n 1 docker rmi --force
> rm -f out/image-id
.PHONY: clean-docker

# Tests - re-run if any Go files have changes since tmp/.tests-passed.sentinel last touched.
tmp/.tests-passed.sentinel: $(shell find . -type f -iname "*.go")
> mkdir -p $(@D)
> go test ./...
> touch $@

# Lint - re-run if the tests have been re-run (and so, by proxy, whenever the source files have changed).
tmp/.linted.sentinel: tmp/.tests-passed.sentinel
> mkdir -p $(@D)
> golangci-lint run
> go vet ./...
> touch $@

# Docker image - re-build if the lint output is re-run.
out/image-id: tmp/.linted.sentinel
> mkdir -p $(@D)
> image_id="$(IMAGE_NAME):$(shell uuidgen)"
> DOCKER_BUILDKIT=1 docker build --tag="$${image_id}" .
> echo "$${image_id}" > out/image-id

# CLI binary - re-build if the lint output is re-run.
mservctl/mservctl: tmp/.linted.sentinel
> cd mservctl
> go build -mod=vendor
