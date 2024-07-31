# Run the app with air for live reload support
dev:
	air -c .air.toml

tidy:
	go mod tidy

lint:	
	golangci-lint run

test:
	go test -race -v -count=1 ./... -cover

build:
	@echo "Linting codebase..."
	golangci-lint run
	@echo "Building binary..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./cmd/app/main ./cmd/app/main.go
	@echo "Build completed successfully..."

run:
	./cmd/app/main

# Metrics service docker compose
metrics\:up:
	docker compose -f ./docker/compose.yml up -d

metrics\:down:
	docker compose -f ./docker/compose.yml down

# Make sure to install migrate cli binary on the system before running this command. https://github.com/golang-migrate/migrate
migrate\:create:
	@echo "Create Migration"
	@read -p "Enter migration name: " migration_name; \
	migrate create -ext sql -dir database/migrations -seq $$migration_name
	@echo "Migration created successfully, fill in the schema in the generated file."
