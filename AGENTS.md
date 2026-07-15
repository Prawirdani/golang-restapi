# AGENTS.md

## Project Overview

Go RESTful API template implementing a clean architecture (Handler -> Service -> Repository) with JWT authentication, PostgreSQL, Redis Streams message queue, Cloudflare R2 storage, and Prometheus monitoring.

- **Module**: `github.com/prawirdani/golang-restapi`
- **Go version**: 1.26.5
- **Router**: go-chi/chi v5
- **Database**: PostgreSQL via pgx v5 (no ORM, raw SQL)
- **Cache/MQ**: Redis (Streams for async email events)
- **Storage**: Cloudflare R2 (S3-compatible)
- **Testing**: testify + mockery (unit tests only, no integration/e2e)

## Key Commands

```bash
make dev            # API server with hot-reload (Air)
make dev:worker     # Background worker with hot-reload
make build          # Build binary (CGO_ENABLED=0, linux)
make run            # Run compiled binary
make test           # go test -race -v -count=1 ./... -cover
make lint           # golangci-lint run
make tidy           # go mod tidy
make migration:up   # Apply goose migrations
make migration:down # Rollback migrations
```

## Project Structure

```
cmd/api/             # API entrypoint (main.go, server.go, container.go)
cmd/worker/          # Background worker entrypoint
config/              # Env-based config structs (App, Postgres, Redis, Auth, CORS, SMTP, R2)
internal/
  domain/            # Business logic (zero infrastructure imports)
    auth/            # Auth service: JWT, sessions, password recovery, crypto
    user/            # User service: CRUD, profile picture
  infrastructure/
    repository/postgres/  # pgx-based DB implementations
    storage/r2/          # Cloudflare R2 (S3-compatible) storage
    messaging/redis/     # Redis Streams producer/consumer
  transport/http/
    context.go       # Custom Context, Handler adapter, error normalization
    handler/         # HTTP handlers (auth_handler, user_handler)
    middleware/       # Auth, CORS, Gzip, RateLimit, Logger, PanicRecovery, RequestID
  worker/            # Email event consumer (Redis Streams -> SMTP)
pkg/
  log/               # Swappable logger (zerolog/slog adapters)
  mailer/            # SMTP email sender (gomail)
  metrics/           # Prometheus instrumentation
  nullable/          # Generic Nullable[T] for optional DB fields
  strings/           # String utilities
  validator/         # Validation wrapper
migrations/          # Goose SQL migrations
```

## Architecture Principles

- **Clean/Onion architecture**: Domain layer has zero infrastructure imports. Interfaces in domain, implementations in infrastructure.
- **Dependency inversion**: Services depend on interfaces (`user.Repository`, `auth.Repository`, `storage.Storage`, `repository.Transactor`). Implementations in `infrastructure/`.
- **DI Container**: `cmd/api/container.go` manually wires all dependencies (no DI framework).
- **3-layer separation**: Handler (HTTP only) -> Service (business logic) -> Repository (data access).

## Code Conventions

### Naming
- Packages: short, lowercase, single-word (`auth`, `user`, `postgres`, `middleware`)
- Files: snake_case (`access_token.go`, `user_repository.go`, `service_test.go`)
- Constructors: `New<Name>()` (e.g., `NewUserRepository`, `NewService`)
- Errors: `Err` prefix (`ErrNotFound`, `ErrWrongCredentials`)
- Constants: PascalCase (`ImageStoragePath`, `AccessTokenCookie`)
- Aliased imports: `httpx` for transport/http, `redisstream` for messaging/redis, `strs` for strings

### Handler Pattern
Handlers use `httpx.Func` signature `func(c *httpx.Context) error` wrapped by `httpx.Handler()` which auto-normalizes errors to JSON responses. **Never manually write error responses** -- just return the error.

```go
func (h *AuthHandler) Login(c *httpx.Context) error {
    // ... logic ...
    return c.JSON(http.StatusOK, httpx.Body{Data: result, Message: "success"})
}
```

### Error Handling
- Typed domain errors: `domain.Error` with `ErrorKind` enum (`KindValidation`, `KindNotFound`, `KindConflict`, `KindUnauthorized`, `KindForbidden`)
- Immutable: `WithDetails()` and `SetMessage()` return copies
- `errors.Is()` support via custom `Is()` method (kind + code comparison)
- `httpx.NormalizeError()` maps domain errors -> HTTP status codes automatically

### Transaction Pattern
Services use `s.transactor.Transact(ctx, func(ctx context.Context) error { ... })`. The tx is injected into context; repositories detect it via `db.GetConn(ctx)` and use the same connection. Repositories add `FOR UPDATE` when they detect a tx connection.

```go
err := s.transactor.Transact(ctx, func(ctx context.Context) error {
    // repositories called here automatically participate in the tx
    return nil
})
```

### Testing
- Table-driven tests with `t.Run` subtests
- `setupTestFixture(t)` creates all mocks, wires the service, registers `t.Cleanup()` for assertion verification
- Mocks generated via `mockery` (config: `.mockery.yml`)
  - `internal/domain/<entity>/mocks/` -- Entity related mocks, eg. Repository, Event-Bus, and other Deps mocks.
  - `internal/testing/mocks/` -- Global reusable mocks.
- Run tests: `make test`
- Generate mocks: `mockery` (reads `.mockery.yml`)

### Logging
- Structured, context-aware via `pkg/log/`
- Logger set at startup: `log.SetLogger(log.NewZerologAdapter(cfg.IsProduction()))`
- Request-scoped fields (request_id, auth uid/sid) propagated via context
- Dev: zerolog Debug level. Prod: zerolog Info level.

## API Routes
API routes registry at `setupRoutes` under `./cmd/api/server.go` 

## Environment

Configuration via `.env` file (see `.env.example`). Config struct parsed at startup by `config.LoadConfig()`. Key env groups: App, Postgres, Redis, CORS, Auth (JWT), SMTP, R2.

## Adding New Features

1. Define domain model and service interface in `internal/domain/<entity>/`, don't forget to add GODOCS on top every new package/entity.
2. Implement repository interface in `internal/infrastructure/repository/postgres/`
3. Add service implementation in `internal/domain/<entity>/`
4. Add HTTP handler in `internal/transport/http/handler/`
5. Wire dependencies in `cmd/api/container.go`
6. Register routes in `cmd/api/server.go` (`setupHandlers`)
7. Add migration in `migrations/` via `make migration:create`
8. Add mocks: update `.mockery.yml` if new interfaces, run `mockery`
9. Write unit tests alongside source files
10. Run `make lint && make test` to verify
