DB_HOST=$(shell python3 -c "import yaml; print(yaml.safe_load(open('./config/config.yml'))['db']['Host'])")
DB_PORT=$(shell python3 -c "import yaml; print(yaml.safe_load(open('./config/config.yml'))['db']['Port'])")
DB_USER=$(shell python3 -c "import yaml; print(yaml.safe_load(open('./config/config.yml'))['db']['Username'])")
DB_PASSWORD=$(shell python3 -c "import yaml; print(yaml.safe_load(open('./config/config.yml'))['db']['Password'])")
DB_NAME=$(shell python3 -c "import yaml; print(yaml.safe_load(open('./config/config.yml'))['db']['Name'])")

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
	goose -dir database/migrations create $$migration_name sql
	@echo "Migration created successfully. Please update the generated file to define the schema changes before running the migration."

# Metrics service docker compose
metrics\:start:
	@echo "Starting metrics service..."
	@if ! docker compose -f ./docker/compose.yml start; then \
		echo "Make sure you already build the metrics service"; \
		exit 1; \
	fi

metrics\:stop:
	docker compose -f ./docker/compose.yml stop

metrics\:build:
	docker compose -f ./docker/compose.yml up -d --build

metrics\:delete:
	docker compose -f ./docker/compose.yml down

