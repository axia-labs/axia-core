.PHONY: build test lint clean install all

all: clean lint test build

build:
	go build -o bin/trust cmd/trust/main.go

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -f bin/trust

install: build
	cp bin/trust /usr/local/bin/

run: build
	./bin/trust

# Development helpers
dev-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go get github.com/sirupsen/logrus
	go get github.com/looplab/fsm
	go get golang.org/x/crypto/sha3

# Environment variables for running the application
export TATUM_API_KEY ?= your-api-key-here
export AXIA_SECRET_KEY ?= $(shell [ -f .env ] && grep AXIA_SECRET_KEY .env | cut -d '=' -f2)

.PHONY: upload-graph
upload-graph:
	@$(MAKE) build
	./bin/trust ipfs upload

.PHONY: get-graph
get-graph:
	@$(MAKE) build
	./bin/trust ipfs get $(IPFS_ID)

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build      - Build the trust binary"
	@echo "  test       - Run tests"
	@echo "  lint       - Run linter"
	@echo "  clean      - Remove built binary"
	@echo "  install    - Install binary to /usr/local/bin"
	@echo "  all        - Run clean, lint, test, and build"
	@echo "  dev-deps   - Install development dependencies"

# Add new target for generating a secret key
.PHONY: generate-key
generate-key:
	@openssl rand -hex 32 > .env
	@echo "Generated new secret key in .env file" 