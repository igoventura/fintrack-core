# Load environment variables from .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: migrate rollback new-migration tidy test help compose

help:
	@echo "Available commands:"
	@echo "  make migrate        - Run database migrations using tern"
	@echo "  make rollback       - Rollback the last migration using tern"
	@echo "  make new-migration  - Create a new migration file (usage: make new-migration name=xxx)"
	@echo "  make tidy           - Run go mod tidy"
	@echo "  make test           - Run all tests"
	@echo "  make compose        - Run docker compose"
	@echo "  make swagger        - Generate OpenAPI documentation using swag"

compose:
	@echo "Running docker compose..."
	@docker compose up -d --build

migrate:
	@echo "Running migrations..."
	@go run github.com/jackc/tern/v2 migrate -m migrations/ -c migrations/tern.conf

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

swagger:
	@echo "Generating Swagger documentation..."
	@go run github.com/swaggo/swag/cmd/swag init -g cmd/api/main.go -o docs --ot yaml

install-hooks:
	@echo "Installing git hooks..."
	@chmod +x scripts/pre-commit
	@ln -sf ../../scripts/pre-commit .git/hooks/pre-commit
	@echo "Hooks installed successfully!"
