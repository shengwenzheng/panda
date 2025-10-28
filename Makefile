# Makefile to run project directly (no build, no tests)
SHELL := /bin/bash

CMD := cmd/api/main.go
CONFIG ?= config/config_test.yaml

GORUN := go run
GOBUILD := go build
GOFMT := gofmt -s -w

.PHONY: all run build fmt vet clean

all: run

run:
	@echo "Running $(CMD) with config: $(CONFIG)"
	$(GORUN) $(CMD) --config-file=$(CONFIG)

build:
	@echo "Build binary..."
	@mkdir -p bin
	$(GOBUILD) -o bin/staker-api ./cmd/api
	@echo "Built bin/staker-api"

start: build
	@echo "Starting bin/staker-api with config: $(CONFIG)"
	./bin/staker-api --config-file=$(CONFIG)


fmt:
	@echo "Formatting..."
	$(GOFMT) .

vet:
	@echo "Running go vet..."
	go vet ./...

clean:
	@echo "Nothing to clean (no binary produced)"