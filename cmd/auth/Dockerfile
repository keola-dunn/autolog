# Stage 1: Builder
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application, making sure it's a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-api ./cmd/auth/main.go

# Stage 2: Runner
FROM alpine:latest AS runner

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/auth-api .

# Expose the port your application listens on (if applicable)
EXPOSE 8080

# Command to run the application
CMD ["./auth-api"]