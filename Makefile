# Makefile to run project directly (no build, no tests)
SHELL := /bin/bash

CMD := cmd/api/main.go
CONFIG ?= config/config_test.yaml

GORUN := go run
GOFMT := gofmt -s -w

.PHONY: all run fmt vet clean

all: run

run:
    @echo "Running $(CMD) with config: $(CONFIG)"
    $(GORUN) $(CMD) --config-file=$(CONFIG)

fmt:
	@echo "Formatting..."
	$(GOFMT) .

vet:
	@echo "Running go vet..."
	go vet ./...

clean:
	@echo "Nothing to clean (no binary produced)"