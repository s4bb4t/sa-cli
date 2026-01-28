BINARY_NAME := sac
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -ldflags "-s -w \
	-X github.com/s4bb4t/sa-cli/internal/version.Version=$(VERSION) \
	-X github.com/s4bb4t/sa-cli/internal/version.GitCommit=$(GIT_COMMIT) \
	-X github.com/s4bb4t/sa-cli/internal/version.BuildDate=$(BUILD_DATE)"

.PHONY: build install clean test lint

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/sac

install:
	go install $(LDFLAGS) ./cmd/sac

clean:
	rm -rf bin/

test:
	go test -v ./...

lint:
	golangci-lint run
