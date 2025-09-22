# Stage 1: Builder
FROM golang:1.23 AS builder

WORKDIR /app

# Copy application code
COPY . .

# Build the Go application, making sure it's a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o auth-api ./cmd/auth/main.go

# Stage 2: Runner
FROM alpine:latest AS runner

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/auth-api .
COPY --from=builder /app/private-key.pem .
COPY --from=builder /app/public-key.pem .

# Expose the port your application listens on (if applicable)
EXPOSE 8080

# Command to run the application
CMD ["./auth-api"]