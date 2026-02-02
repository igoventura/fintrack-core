# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy project files
COPY . .

# Build the application
RUN go build -o /fintrack-api ./cmd/api

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /fintrack-api .
# Copy documentation and migrations as they are needed at runtime
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/migrations ./migrations

# Expose the API port
EXPOSE 8080

# Run the application
CMD ["./fintrack-api"]
