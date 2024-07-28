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

# Lint Code
lint:	
	golangci-lint run

tidy:
	go mod tidy

test:
	go test -v -count=1 ./... -cover

docker\:prod:
	DOCKERFILE=Dockerfile  docker compose up -d --build

# Make sure to install migrate cli binary on the system before running this command. https://github.com/golang-migrate/migrate
migrate\:create:
	@echo "Create Migration"
	@read -p "Enter migration name: " migration_name; \
	migrate create -ext sql -dir database/migrations -seq $$migration_name
	@echo "Migration created successfully, fill in the schema in the generated file."
