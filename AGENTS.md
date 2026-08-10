# AGENTS.md

Go RESTful API template. Clean architecture (Handler -> Service -> Repository) with JWT auth, PostgreSQL, Redis Streams, Cloudflare R2, and Prometheus.

## Stack

- **Module**: `github.com/prawirdani/golang-restapi` (Go 1.26.5)
- **Router**: go-chi/chi v5
- **DB**: PostgreSQL via pgx v5 — raw SQL, no ORM (goose migrations)
- **MQ/Cache**: Redis Streams (async email events)
- **Storage**: Cloudflare R2 (S3-compatible)
- **Testing**: testify + mockery (unit tests only)

## Commands

```bash
make dev            # API server, hot-reload (Air)
make dev:worker     # Worker, hot-reload
make build          # Build binary (CGO_ENABLED=0, linux)
make test           # go test -race -count=1 ./... -cover
make lint           # golangci-lint run
make migration:create  # Scaffold a goose migration
make migration:up      # Apply migrations
mockery             # Regenerate mocks (reads .mockery.yml)
```

## Layout

```
cmd/api/                 # API entrypoint: main.go, server.go (routes), container.go (DI)
cmd/worker/              # Background worker entrypoint
config/                  # Env-based config (App, Postgres, Redis, Auth, CORS, SMTP, R2)
internal/
  domain/                # Business logic — ZERO infrastructure imports
    auth/                # JWT, sessions, password recovery, crypto
    user/                # CRUD, profile picture
  infrastructure/
    repository/postgres/ # pgx repository implementations
    storage/r2/          # R2 storage
    messaging/redis/     # Redis Streams producer/consumer
  transport/http/        # context.go (Context, Handler, error normalization), handler/, middleware/
  worker/                # Email event consumer (Redis -> SMTP)
pkg/                     # log, mailer, metrics, nullable, strings, validator
migrations/              # Goose SQL migrations
```

## Architecture

- **Onion/Clean**: interfaces live in `domain/`, implementations in `infrastructure/`. Domain imports no infrastructure.
- **Dependency inversion**: services depend on interfaces (`user.Repository`, `auth.Repository`, `storage.Storage`, `repository.Transactor`).
- **DI**: `cmd/api/container.go` wires everything manually — no framework.
- Handler (HTTP only) -> Service (business logic) -> Repository (data access).

## Conventions

**Naming**
- Packages: short, lowercase, single-word (`auth`, `postgres`, `middleware`)
- Files: snake_case (`user_repository.go`, `service_test.go`)
- Constructors: `New<Name>()`; Errors: `Err` prefix; Constants: PascalCase
- Import aliases: `httpx` (transport/http), `redisstream` (messaging/redis), `strs` (pkg/strings)
- Add package-level godoc to every new package/entity.

**Handlers** — signature `func(c *httpx.Context) error`, wrapped by `httpx.Handler()`. Return errors; never write error responses manually.

```go
func (h *AuthHandler) Login(c *httpx.Context) error {
    return c.JSON(&httpx.Body{Data: result, Message: "success"})
}
```

**Errors** — `domain.Error` with `ErrorKind` (`KindValidation`, `KindNotFound`, `KindConflict`, `KindUnauthorized`, `KindForbidden`). Immutable (`WithDetails`/`SetMessage` return copies), supports `errors.Is`. `httpx.NormalizeError` maps kinds to HTTP status.

**Transactions** — wrap multi-step writes in `s.transactor.Transact`. Repositories detect the tx via `db.GetConn(ctx)` and reuse the connection (adding `FOR UPDATE`).

```go
err := s.transactor.Transact(ctx, func(ctx context.Context) error {
    // repos here automatically join the tx
    return nil
})
```

**Testing** — table-driven with `t.Run` subtests. `setupTestFixture(t)` wires mocks + service and registers `t.Cleanup()`. Mocks: entity-scoped in `internal/domain/<entity>/mocks/`, reusable in `internal/testing/mocks/`.

**Logging** — structured/context-aware via `pkg/log`. Set at startup: `log.SetLogger(log.NewZerologAdapter(cfg.IsProduction()))`. Request-scoped fields (request_id, uid/sid) flow through context. Debug in dev, Info in prod.

**Config** — `.env` (see `.env.example`), loaded once via `config.LoadConfig()`.

## Custom Skills

Project-specific skills in `.agents/skills/` encode this repo's style, architecture, and guidelines. Load them (they auto-trigger) whenever working on the related layer:

- **gorest-architecture** — layering, dependency direction, interfaces, DI wiring, aliases, naming
- **gorest-errors** — domain error constructors/kinds, immutable copies, repo error translation
- **gorest-handlers** — httpx handler signature, BindValidate, response envelope, cookies, multipart
- **gorest-repositories** — pgx builders, pgxscan, tx-aware GetConn/FOR UPDATE, error mapping
- **gorest-services** — Transact, post-commit side effects, nullable+Validate, async cleanup
- **gorest-testing** — mockery placement, setupTestFixture, transactor expectations, subtests

These supersede generic samber guidance where they overlap.

## Adding a Feature

1. Define model + service interface in `internal/domain/<entity>/` (with godoc).
2. Implement repository in `internal/infrastructure/repository/postgres/`.
3. Add service implementation in `internal/domain/<entity>/`.
4. Add handler in `internal/transport/http/handler/`.
5. Wire in `cmd/api/container.go`; register routes in `cmd/api/server.go` (`setupHandlers`).
6. Add migration via `make migration:create`.
7. Update `.mockery.yml` for new interfaces, run `mockery`.
8. Write unit tests alongside source; verify with `make lint && make test`.
