## Golang REST API
Personal Go RESTful API template with common 3 Layered architecture with following layers:
- `Handler/Controller/Delivery`: Responsible for handling incoming HTTP request, parsing request, validating request, calling service layer and sending response.
- `Service/Usecase`: Business logic layer, responsible for handling business logic, calling repository layer and returning response to handler.
- `Repository/Store`: Responsible for handling database operation, query, insert, update, delete etc.

### Lib/Tool :
- HTTP router: [chi](https://github.com/go-chi/chi)
- Struct validator: [validator](https://github.com/go-playground/validator)
- Unique Identifier: [uuid](https://github.com/google/uuid)
- Postgres driver & pooling: [pgx](https://github.com/jackc/pgx)
- Postgres struct scanner: [scanny](https://github.com/georgysavva/scany)
- Cloudflare R2 Storage Using [AWS-sdk-v2](https://github.com/aws/aws-sdk-go-v2)
- Message Queueing: [RabbitMQ](https://github.com/rabbitmq/amqp091-go)
- SMTP Mailing: [gomail.v2](https://pkg.go.dev/gopkg.in/gomail.v2)
- ENV Loader: [godotenv](https://github.com/joho/godotenv)
- Logger: std `log/slog` & [zerolog](https://github.com/rs/zerolog) (Swappable)
- Auth: [jwt](https://github.com/golang-jwt/jwt)
- Testing: [testify](https://github.com/stretchr/testify)
- Metrics & Instrumentation:
    - Prometheus: [prometheus](https://github.com/prometheus/client_golang)
    - Grafana: [grafana](https://grafana.com)
- CLI:
    - Database Migration Tool: [goose](https://github.com/pressly/goose)
    - Development live reloading: [air](https://github.com/cosmtrek/air)
    - Linters: [golangci-lint](https://github.com/golangci/golangci-lint)

