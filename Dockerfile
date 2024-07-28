FROM golang:1.21.1-alpine AS builder

WORKDIR /app

RUN apk update && apk add --no-cache make

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# RUN make test

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

# Final Minimal Image
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main ./main

# Copy config files
COPY --from=builder /app/config ./config

ENTRYPOINT ["./main"]
