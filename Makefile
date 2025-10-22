# Load .env variables
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Run the app with air for live reload support
dev:
	@air -c .air.toml

tidy:
	@go mod tidy

lint:	
	@golangci-lint run

test:
	@go test -race -v -count=1 ./... -cover

build:
	@echo "Building binary..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./cmd/app/main ./cmd/app/main.go
	@echo "Build completed successfully..."

run:
	./cmd/app/main

# Makesure you have goose binary installed
migration\:status:
	@goose -dir database/migrations postgres "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" status

migration\:up:
	@goose -dir database/migrations postgres "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" up

migration\:down:
	@goose -dir database/migrations postgres "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" down

migration\:create:
	@echo "Create Migration"
	@read -p "Enter migration name: " migration_name; \
	goose -s -dir database/migrations create $$migration_name sql
	@echo "Migration created successfully, fill in the schema in the generated file."

# Metrics service docker compose
metrics\:start:
	@echo "Starting metrics service..."
	@if ! docker compose -f ./deployment/docker/compose.yml start; then \
		echo "Make sure you already build the metrics service"; \
		exit 1; \
	fi

metrics\:stop:
	docker compose -f ./deployment/docker/compose.yml stop

metrics\:build:
	docker compose -f ./deployment/docker/compose.yml up -d --build

metrics\:delete:
	docker compose -f ./deployment/docker/compose.yml down

