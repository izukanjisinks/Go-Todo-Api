# syntax=docker/dockerfile:1

# Build stage
FROM golang:alpine AS builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Tidy modules and build
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o /app/todo-api ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy binary from builder stage
COPY --from=builder /app/todo-api /app/todo-api

# Expose port (adjust to match your SERVER_PORT in config)
EXPOSE 8080

# Run the application
CMD ["/app/todo-api"]