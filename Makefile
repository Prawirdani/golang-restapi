# Load .env variables
ifneq (,$(wildcard ./.env))
    include .env
    export
endif


# Run the api server
dev:
	@air -c .air.toml

# Run the mqworker
dev\:mqworker:
	@air -c .air.mqworker.toml

tidy:
	@go mod tidy

lint:	
	@golangci-lint run

test:
	@go test -race -v -count=1 ./... -cover

build:
	@echo "Building binary..."
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o ./bin/api ./cmd/api/main.go
	@echo "Build completed successfully..."

run:
	./bin/api

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


compose:
	@echo "🚀 Composing....."
	docker compose  up -d --build

compose-down:
	@echo "🧹 Stopping docker compose"
	docker compose down

compose-logs:
	docker compose logs -f api

