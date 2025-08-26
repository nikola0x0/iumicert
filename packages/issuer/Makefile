.PHONY: build run dev serve clean init generate batch

# Build the binary
build:
	go build -o micert cmd/*.go

# Run the server in development mode
dev:
	go run cmd/*.go serve --port 8080 --cors

# Alternative serve command
serve: dev

# Initialize repository
init:
	go run cmd/*.go init "IU-VNUHCM"

# Generate test data
generate:
	go run cmd/*.go generate-data

# Batch process all terms
batch:
	go run cmd/*.go batch-process

# Build and run
run: build
	./micert serve --port 8080 --cors

# Clean build artifacts
clean:
	rm -f micert

# Install dependencies
deps:
	go mod tidy
	go mod download

# Quick setup (init + generate + batch)
setup: init generate batch

# Development workflow
start: clean build run