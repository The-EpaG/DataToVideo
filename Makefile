# Simple Makefile

# Build the application
all: build

build:
	@echo "Building..."

	@mkdir -p bin
	@go build -o bin/DataToVideo cmd/api/main.go

	@echo "Done."

# Run the application
run:
	@go run cmd/api/main.go