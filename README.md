## Golang REST API
Simple, structured, easy to use and new commer friendly golang REST API project.

### Architecture
Common Layered architecture with following layers:
- `Handler/Controller/Delivery`: Responsible for handling incoming HTTP request, parsing request, validating request, calling service layer and sending response.
- `Service/Usecase`: Business logic layer, responsible for handling business logic, calling repository layer and returning response to handler.
- `Repository/Store`: Responsible for handling database operation, query, insert, update, delete etc.

### Lib/Tool :
- HTTP router: [chi](https://github.com/go-chi/chi)
- Struct validator: [validator](https://github.com/go-playground/validator)
- Unique Identifier: [uuid](https://github.com/google/uuid)
- Postgres driver & pooling: [pgx](https://github.com/jackc/pgx)
- Postgres struct scanner: [scanny](https://github.com/georgysavva/scany)
- Config parser: [viper](https://github.com/spf13/viper)
- Logger: [zerolog](https://github.com/rs/zerolog) 
- Auth: [jwt](https://github.com/golang-jwt/jwt)
- Testing: [testify](https://github.com/stretchr/testify)
- Metrics & Instrumentation:
    - Prometheus: [prometheus](https://github.com/prometheus/client_golang)
    - Grafana: [grafana](https://grafana.com)
- CLI:
    - Database Migration Tool: [goose](https://github.com/pressly/goose)
    - Development live reloading: [air](https://github.com/cosmtrek/air)
    - Linters: [golangci-lint](https://github.com/golangci/golangci-lint)


### Configuration
Application configuration is in `config/config.example.yml`, rename it to `config.yml` before you proceed.

### Database Migration
Migration is handled by `goose` cli tool. make sure you have installed goose binary in your system. All migrations file is in `database/migrations` directory. 
##### Create Migration
```bash
# To Create new db migration.
make migrate:create
```
##### Run Migration
```bash
# To Execute/Run the migration file.
goose -dir database/migrations postgres "host=localhost port=5432 user=<your_username> password=<your_password> dbname=<your_db_name> sslmode=disable" up
```
Please refer to [goose](https://github.com/pressly/goose) for more detail documentation & instruction.

### Docker
The metrics and instrumentation are configured to run in Docker, while the main application service runs on the host machine.
To build and run the Metrics service using Docker Compose, refer to the commands in the Makefile.
