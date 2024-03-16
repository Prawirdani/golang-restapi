# Run the app with air for live reload support
dev:
	air -c .air.toml

# Build binary
build:
	@echo "Linting codebase..."
	golangci-lint run
	@echo "Building binary..."
	go build -o ./cmd/bin ./cmd/main.go
	@echo "Build completed successfully..."

# Run binary
run:
	@echo "Running binary..."
	./cmd/bin

# Lint Code
lint:	
	golangci-lint run

tidy:
	go mod tidy