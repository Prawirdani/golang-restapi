# Run the app with air for live reload support
dev:
	air -c .air.toml

# Build binary
build:
	@echo "Linting codebase..."
	golangci-lint run
	@echo "Building binary..."
	go build -o ./cmd/app/bin ./cmd/app/main.go
	@echo "Build completed successfully..."

# Run binary
run:
	@echo "Running binary..."
	./cmd/app/bin

# Lint Code
lint:	
	golangci-lint run

tidy:
	go mod tidy

test:
	go test -v -count=1 ./... -cover
