SHELL := bash

# Default - top level rule is what gets run when you run just `make` without specifying a goal/target.
.DEFAULT_GOAL := help

.DELETE_ON_ERROR:
.ONESHELL:
.SHELLFLAGS := -euo pipefail -c

MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --warn-undefined-variables

export TYK_VERSION := v5.6.1

ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later.)
endif
.RECIPEPREFIX = >

image_repository ?= tykio/mserv

# Adjust the width of the first column by changing the '16' value in the printf pattern.
help:
> @grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
  | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'
.PHONY: help

check-swagger:
> which swagger || (go get -u github.com/go-swagger/go-swagger/cmd/swagger)

swagger: check-swagger
> swagger generate spec -o ./doc/swagger.yaml --scan-models -x mservclient
.PHONY: swagger

serve-swagger: check-swagger
> swagger serve -F=swagger ./doc/swagger.yaml

swagger-client: check-swagger
> mkdir -p ./mservclient
> swagger generate client -f ./doc/swagger.yaml -t ./mservclient

clean: ## Clean up the temp and output directories, and any built binaries. This will cause everything to get rebuilt.
> rm -rf ./bin
> rm -rf ./integration/outputs
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
> CGO_ENABLED=0 hack/bin/golangci-lint run

hack/bin/golangci-lint:
> curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
> | sh -s -- -b $(shell pwd)/hack/bin

docker: Dockerfile ## Builds mserv docker image.
> mkdir -p $(@D)
> image_id="$(image_repository):$(shell uuidgen)"
> DOCKER_BUILDKIT=1 docker build --tag="$${image_id}" .

build: mserv mservctl ## Build server and client binary.
.PHONY: build

mserv:
> CGO_ENABLED=0 go build -o bin/mserv
.PHONY: mserv

mservctl:
> cd mservctl
> CGO_ENABLED=0 go build -o ../bin/mservctl
.PHONY: mservctl

start: ## Start runs development environment with mserv and mongo in docker compose.
> docker compose up -d

stop: ## Stop runs development environment with mserv and mongo in docker compose.
> docker compose stop

# Builds multiple Go plugins and moves them into local Tyk instance.
plugins:
> @for plugin in plugin_1.go plugin_2.go; do \
>	docker compose run --rm tyk-plugin-compiler $$plugin _$$(date +%s); \
> done
.PHONY: plugins

bundles:
> docker compose run --rm --workdir /plugin-source --entrypoint "/opt/tyk-gateway/tyk bundle build -y -o bundle.zip" tyk-gateway
.PHONY: bundles

integration: ## Runs integration test for mserv and mservctl it needs services running.
> cd integration
> venom run integration.yaml -vvv --output-dir outputs
> cd ..
.PHONY: integration
