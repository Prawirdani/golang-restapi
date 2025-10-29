# =========================
# Stage 1: Build
# =========================
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache make 

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (for caching dependencies)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build using Makefile
RUN make build

# =========================
# Stage 2: Minimal runtime image
# =========================
FROM alpine:3.22

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/bin/api .
COPY --from=builder /app/templates ./templates

# Expose port (if applicable)
EXPOSE 8080

# Run the binary
CMD ["./api"]
