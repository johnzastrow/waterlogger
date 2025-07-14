# Waterlogger Makefile

# Variables
APP_NAME := waterlogger
VERSION := 1.0.0
BUILD_DIR := build
DIST_DIR := dist
MAIN_FILE := cmd/waterlogger/main.go

# Go build flags
LDFLAGS := -ldflags "-X main.version=$(VERSION) -s -w"
BUILD_FLAGS := -trimpath

# Default target
.PHONY: all
all: build

# Build for current platform
.PHONY: build
build:
	go build $(BUILD_FLAGS) $(LDFLAGS) -o $(APP_NAME) $(MAIN_FILE)

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR) $(DIST_DIR)
	rm -f $(APP_NAME) $(APP_NAME).exe waterlogger-mac
	rm -f *.db *.log

# Run tests
.PHONY: test
test:
	go test ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run tests with race detection
.PHONY: test-race
test-race:
	go test -race ./...

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	go vet ./...
	golint ./...

# Tidy dependencies
.PHONY: tidy
tidy:
	go mod tidy

# Download dependencies
.PHONY: deps
deps:
	go mod download

# Build for all platforms
.PHONY: build-all
build-all: build-linux build-windows build-mac

# Build for Linux
.PHONY: build-linux
build-linux:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(MAIN_FILE)

# Build for Windows
.PHONY: build-windows
build-windows:
	mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(MAIN_FILE)

# Build for macOS
.PHONY: build-mac
build-mac:
	mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(MAIN_FILE)

# Create distribution packages
.PHONY: dist
dist: build-all
	mkdir -p $(DIST_DIR)
	
	# Linux package
	mkdir -p $(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-amd64
	cp $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-amd64/$(APP_NAME)
	cp config.yaml $(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-amd64/config.example.yaml
	cp README.md $(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-amd64/
	cp LICENSE $(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-amd64/
	cd $(DIST_DIR) && tar -czf $(APP_NAME)-$(VERSION)-linux-amd64.tar.gz $(APP_NAME)-$(VERSION)-linux-amd64/
	
	# Windows package
	mkdir -p $(DIST_DIR)/$(APP_NAME)-$(VERSION)-windows-amd64
	cp $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(DIST_DIR)/$(APP_NAME)-$(VERSION)-windows-amd64/$(APP_NAME).exe
	cp config.yaml $(DIST_DIR)/$(APP_NAME)-$(VERSION)-windows-amd64/config.example.yaml
	cp README.md $(DIST_DIR)/$(APP_NAME)-$(VERSION)-windows-amd64/
	cp LICENSE $(DIST_DIR)/$(APP_NAME)-$(VERSION)-windows-amd64/
	cd $(DIST_DIR) && zip -r $(APP_NAME)-$(VERSION)-windows-amd64.zip $(APP_NAME)-$(VERSION)-windows-amd64/
	
	# macOS package
	mkdir -p $(DIST_DIR)/$(APP_NAME)-$(VERSION)-darwin-amd64
	cp $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(DIST_DIR)/$(APP_NAME)-$(VERSION)-darwin-amd64/$(APP_NAME)
	cp config.yaml $(DIST_DIR)/$(APP_NAME)-$(VERSION)-darwin-amd64/config.example.yaml
	cp README.md $(DIST_DIR)/$(APP_NAME)-$(VERSION)-darwin-amd64/
	cp LICENSE $(DIST_DIR)/$(APP_NAME)-$(VERSION)-darwin-amd64/
	cd $(DIST_DIR) && tar -czf $(APP_NAME)-$(VERSION)-darwin-amd64.tar.gz $(APP_NAME)-$(VERSION)-darwin-amd64/

# Run the application
.PHONY: run
run:
	./$(APP_NAME)

# Run the application with hot reload (requires air)
.PHONY: dev
dev:
	air

# Install development dependencies
.PHONY: install-dev
install-dev:
	go install github.com/cosmtrek/air@latest
	go install golang.org/x/lint/golint@latest

# Database operations
.PHONY: db-reset
db-reset:
	rm -f waterlogger.db
	./$(APP_NAME) &
	sleep 2
	pkill -f $(APP_NAME)

# Generate documentation
.PHONY: docs
docs:
	go doc ./...

# Check for security vulnerabilities
.PHONY: security
security:
	go list -json -m all | nancy sleuth

# Performance benchmarks
.PHONY: bench
bench:
	go test -bench=. ./...

# Profile the application
.PHONY: profile
profile:
	go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=. ./...

# Install the application
.PHONY: install
install:
	go install $(LDFLAGS) $(MAIN_FILE)

# Create example configuration
.PHONY: config
config:
	cp config.yaml config.example.yaml

# Docker build (if Dockerfile exists)
.PHONY: docker-build
docker-build:
	docker build -t $(APP_NAME):$(VERSION) .

# Docker run (if Dockerfile exists)
.PHONY: docker-run
docker-run:
	docker run -p 2341:2341 $(APP_NAME):$(VERSION)

# Help target
.PHONY: help
help:
	@echo "Waterlogger Build System"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all           Build the application for current platform"
	@echo "  build         Build the application for current platform"
	@echo "  build-all     Build for all platforms (Linux, Windows, macOS)"
	@echo "  build-linux   Build for Linux"
	@echo "  build-windows Build for Windows"
	@echo "  build-mac     Build for macOS"
	@echo "  clean         Clean build artifacts"
	@echo "  test          Run tests"
	@echo "  test-coverage Run tests with coverage report"
	@echo "  test-race     Run tests with race detection"
	@echo "  fmt           Format code"
	@echo "  lint          Run linting tools"
	@echo "  tidy          Tidy dependencies"
	@echo "  deps          Download dependencies"
	@echo "  dist          Create distribution packages"
	@echo "  run           Run the application"
	@echo "  dev           Run with hot reload (requires air)"
	@echo "  install-dev   Install development dependencies"
	@echo "  db-reset      Reset database"
	@echo "  docs          Generate documentation"
	@echo "  security      Check for security vulnerabilities"
	@echo "  bench         Run performance benchmarks"
	@echo "  profile       Profile the application"
	@echo "  install       Install the application"
	@echo "  config        Create example configuration"
	@echo "  docker-build  Build Docker image"
	@echo "  docker-run    Run Docker container"
	@echo "  help          Show this help message"