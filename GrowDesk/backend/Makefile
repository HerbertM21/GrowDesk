.PHONY: build run test clean docker docker-compose lint tidy

# Default build directory
BUILD_DIR := ./build
# Binary name
BINARY_NAME := backend

# Build the application
build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server/main.go

# Run the application
run:
	go run ./cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)

# Build docker image
docker:
	docker build -t growdesk-backend:latest .

# Run using docker-compose
docker-compose:
	docker-compose up --build

# Run docker-compose in background
docker-compose-detached:
	docker-compose up -d --build

# Stop docker-compose
docker-compose-down:
	docker-compose down

# Run linter
lint:
	go vet ./...

# Update dependencies
tidy:
	go mod tidy

# Initialize a new module and download dependencies
init:
	go mod init github.com/growdesk/growdesk-go
	go mod tidy 