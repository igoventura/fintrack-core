# Load environment variables from .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: migrate rollback new-migration tidy test help

help:
	@echo "Available commands:"
	@echo "  make migrate        - Run database migrations using tern"
	@echo "  make rollback       - Rollback the last migration using tern"
	@echo "  make new-migration  - Create a new migration file (usage: make new-migration name=xxx)"
	@echo "  make tidy           - Run go mod tidy"
	@echo "  make test           - Run all tests"

migrate:
	@echo "Running migrations..."
	@go run github.com/jackc/tern/v2 migrate -m migrations/

rollback:
	@echo "Rolling back last migration..."
	@go run github.com/jackc/tern/v2 rollback -m migrations/

new-migration:
	@if [ -z "$(name)" ]; then echo "Usage: make new-migration name=some_name"; exit 1; fi
	@go run github.com/jackc/tern/v2 new -m migrations/ $(name)

tidy:
	@go mod tidy

test:
	@echo "Running tests..."
	@go test ./...
