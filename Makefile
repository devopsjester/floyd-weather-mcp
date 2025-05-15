# Makefile for Floyd Weather MCP Server

# Variables
BINARY_NAME=floyd-weather-server
MAIN_FILE=main.go
GO=go
GOFMT=gofmt
GOBUILD=$(GO) build
GOTEST=$(GO) test
GOCLEAN=$(GO) clean
GOGET=$(GO) get

# Default city and country (can be overridden from command line)
CITY?=London
COUNTRY?=United Kingdom

# Build the application
build:
	@echo "Building..."
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_FILE)

# Run the application
run: build
	@echo "Running..."
	./$(BINARY_NAME)
	
# Run the weather check for a city
weather: build
	@echo "Getting weather information..."
	./run.sh weather "$(CITY)" "$(COUNTRY)"
	
# Check deployment safety for a city
check-safety: build
	@echo "Checking deployment safety..."
	./run.sh check "$(CITY)" "$(COUNTRY)"
	
# Deploy to a city
deploy: build
	@echo "Attempting deployment..."
	./run.sh deploy "$(CITY)" "$(COUNTRY)"

# Clean the build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	./test.sh

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -w .

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOGET) ./...

# Display help information
help:
	@echo "Floyd Weather MCP Server Makefile"
	@echo "Usage:"
	@echo "  make build    - Build the application"
	@echo "  make run      - Build and run the application"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make test     - Run tests"
	@echo "  make fmt      - Format code"
	@echo "  make deps     - Install dependencies"
	@echo "  make help     - Display this help"

# Default target
.DEFAULT_GOAL := build
