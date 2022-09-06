#! /usr/bin/make -f

# ---------------------------------------------------------
# Project variables.
# ---------------------------------------------------------
PROJECT_NAME := $(shell basename "$(PWD)")
PACKAGE := github.com/unanoc/$(PROJECT_NAME)
VERSION := $(shell git describe --tags 2>/dev/null || git describe --all)
BUILD := $(shell git rev-parse --short HEAD)
DATETIME := $(shell date +"%Y.%m.%d-%H:%M:%S")

# ---------------------------------------------------------
# Use linker flags to provide version/build settings.
# ---------------------------------------------------------
LDFLAGS=-ldflags "-X=$(PACKAGE)/build.Version=$(VERSION) -X=$(PACKAGE)/build.Build=$(BUILD) -X=$(PACKAGE)/build.Date=$(DATETIME)"

# ---------------------------------------------------------
# Go related variables.
# ---------------------------------------------------------
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin

# ---------------------------------------------------------
# Go files.
# ---------------------------------------------------------
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

# ---------------------------------------------------------
# Docker Variables
# ---------------------------------------------------------
DOCKER_COMPOSE_FILE_LOCAL_ENV ?= ./docker-compose.yaml

# ---------------------------------------------------------
# API
# ---------------------------------------------------------
api: go-build-api start-api

start-api:
	@echo "  >  Starting api"
	PROMETHEUS_SUBSYSTEM=api $(GOBIN)/api

go-build-api:
	@echo "  >  Building api binary..."
	GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/api ./cmd/api

# ---------------------------------------------------------
# Parser
# ---------------------------------------------------------
parser: go-build-parser start-parser

start-parser:
	@echo "  >  Starting parser"
	PROMETHEUS_SUBSYSTEM=parser $(GOBIN)/parser

go-build-parser:
	@echo "  >  Building parser binary..."
	GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/parser ./cmd/parser

# ---------------------------------------------------------
# Code Checking
# ---------------------------------------------------------
check: fmt lint unit-test swag

unit-test:
	@echo "  >  Running unit tests"
	go clean -testcache
	GOBIN=$(GOBIN) go test -cover -race -coverprofile=coverage.txt -covermode=atomic -v ./...

fmt:
	@echo "  >  Format all go files"
	GOBIN=$(GOBIN) gofmt -w ${GOFMT_FILES}

lint-install:
ifeq ("$(wildcard bin/golangci-lint)","")
	@echo "  >  Installing golint"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s
endif

lint: lint-install
	@echo "  >  Running golint"
	bin/golangci-lint run --timeout=3m

swag:
	@echo "  >  Generating swagger files"
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g cmd/api/main.go

# ---------------------------------------------------------
# Local Environment
# ---------------------------------------------------------
up:
	@echo "  >  Run local environment"
	docker-compose -f ${DOCKER_COMPOSE_FILE_LOCAL_ENV} up -d