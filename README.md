## Golang REST API
Simple, structured, easy to use and new commer friendly golang REST API project.

### Architecture
[TBD]

### Lib/Tool :
- HTTP router: [chi](https://github.com/go-chi/chi)
- Struct validator: [validator](https://github.com/go-playground/validator)
- Unique Identifier: [uuid](https://github.com/google/uuid)
- Postgres driver & pooling: [pgx](https://github.com/jackc/pgx)
- Config parser: [viper](https://github.com/spf13/viper)
- Logger: [slog](https://pkg.go.dev/golang.org/x/exp/slog) go builtin logger package, this package only available for `go 1.21+`.
- JWT Auth: [jwt](https://github.com/golang-jwt/jwt)
- CLI:
    - Database Migration Tool: [migrate](https://github.com/golang-migrate/migrate)
    - Development live reloading: [air](https://github.com/cosmtrek/air)
    - Linters: [golangci-lint](https://github.com/golangci/golangci-lint)


### Configuration
Application configuration is in `config.example.json`, rename it to `config.json` before you proceed.

### Database Migration
All migrations file is in `database/migrations` directory.
##### Create Migration
```bash
# To Create new db migration.
migrate create -ext sql -dir database/migrations <migration_name>
```
##### Run Migration
```bash
# To Execute/Run the migration file.
migrate -path database/migrations -database "postgresql://<username>:<password>@localhost:5432/<db-name>?sslmode=disable" -verbose up
```
Please refer to [migrate](https://github.com/golang-migrate/migrate) for detail documentation & instruction.

