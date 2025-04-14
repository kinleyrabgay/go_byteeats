# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=api
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_PATH=cmd/api/main.go

# Docker parameters
DOCKER_COMPOSE=docker-compose
DOCKER_BUILD=docker build

# Database parameters
DB_NAME=go_backend
DB_USER=postgres
DB_PASSWORD=postgres
DB_HOST=localhost
DB_PORT=5432
MIGRATION_DIR=migrations

.PHONY: all build clean test coverage deps dev docker-build docker-run migrate migrate-down lint fmt help

all: clean deps build

build:
	@echo "Building..."
	$(GOBUILD) -o bin/$(BINARY_NAME) $(MAIN_PATH)

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f bin/$(BINARY_NAME)
	rm -f bin/$(BINARY_UNIX)

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) verify

dev:
	@echo "Starting development server..."
	air -c .air.toml

docker-build:
	@echo "Building Docker image..."
	$(DOCKER_BUILD) -t $(BINARY_NAME) .

docker-run:
	@echo "Running Docker container..."
	$(DOCKER_COMPOSE) up

migrate:
	@echo "Running database migrations..."
	migrate -path $(MIGRATION_DIR) -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	@echo "Rolling back database migrations..."
	migrate -path $(MIGRATION_DIR) -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $$name

lint:
	@echo "Running linters..."
	golangci-lint run

fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

generate:
	@echo "Generating mocks..."
	mockgen -source=internal/repository/interfaces.go -destination=internal/repository/mocks/repository_mocks.go
	mockgen -source=internal/services/interfaces.go -destination=internal/services/mocks/service_mocks.go

swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/api/main.go -o docs

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_UNIX) $(MAIN_PATH)

# Help command
help:
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make clean         - Clean build files"
	@echo "  make test          - Run tests"
	@echo "  make coverage      - Run tests with coverage"
	@echo "  make deps          - Download dependencies"
	@echo "  make dev           - Start development server"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run Docker container"
	@echo "  make migrate       - Run database migrations"
	@echo "  make migrate-down  - Rollback migrations"
	@echo "  make migrate-create- Create a new migration"
	@echo "  make lint          - Run linters"
	@echo "  make fmt           - Format code"
	@echo "  make generate      - Generate mocks"
	@echo "  make swagger       - Generate Swagger docs"
	@echo "  make build-linux   - Build for Linux"
