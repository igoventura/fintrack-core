# FinTrack Core

FinTrack Core is a Web API providing the foundational financial tracking logic for the FinTrack ecosystem. It is built with Go, following Clean Architecture principles to ensure modularity, testability, and clear separation of concerns.

## Structure

See [PROJECT_STRUCTURE.md](file:///Users/igoventura/Developer/Personal/fintrack-api/PROJECT_STRUCTURE.md) for details on the codebase organization and architectural layers.

## Infrastructure

The project includes a full Docker setup for development.

```bash
# Start Postgres and automatically run migrations
make compose
```

## Automation

A `Makefile` is provided for common development tasks. It automatically loads variables from a `.env` file if it exists.

```bash
# Infrastructure
make compose        # Start DB and run migrations via Docker

# Migrations (Local)
make migrate        # Run migrations using local Go tools
make rollback       # Rollback the last migration
make new-migration name=xxx  # Create a new migration file

# Maintenance
make tidy           # Run go mod tidy
make test           # Run all tests
make swagger        # Regenerate OpenAPI documentation using swag
make install-hooks  # Install git pre-commit hooks
```

## Migrations

This project uses [tern](https://github.com/jackc/tern) for database migrations. Configuration is handled via `migrations/tern.conf` which pulls from environment variables.

## Getting Started

### Prerequisites

- Go 1.25.6 or higher
- PostgreSQL (for local development/runtime)

### Installation

```bash
go get github.com/igoventura/fintrack-api
```

### Usage

The API can be run locally via Docker (recommended) or standalone.

#### Using Docker (Infrastructure included)
```bash
make compose
```
Once running, the API is available at `http://localhost:8080`.

#### Standalone Execution
```bash
# Ensure .env is configured
go run cmd/api/main.go
```

### Documentation (OpenAPI / Scalar)

The API automatically serves interactive documentation generated from code annotations.

- **Redoc UI**: Replaced with **Scalar** for valid, interactive API documentation.
- **Documentation URL**: [http://localhost:8080/docs](http://localhost:8080/docs)
- **OpenAPI Spec**: [http://localhost:8080/swagger.yaml](http://localhost:8080/swagger.yaml)

To update the documentation after changing code annotations, run:
```bash
make swagger
```

### Health Check
```bash
curl http://localhost:8080/health
```

## Authentication

FinTrack Core uses **Supabase Authentication** for secure identity management.
- **Identity Provider**: Supabase Auth (JWT).
- **Validation**: Server-side JWT validation using the `internal/auth` package, which verifies tokens against Supabase's JWKS.
- **Internal Mapping**: Users are linked via a `supabase_id` column in the `users` table.
- **Middleware**: `AuthMiddleware` extracts the Bearer token, validates it, and injects the user context into the request.

## Multi-tenancy

FinTrack Core supports strict multi-tenancy via request headers.
- **Header**: `X-Tenant-ID` (Required)
- **Validation**: `TenantMiddleware` validates the existence of the tenant in the database. Returns `401 Unauthorized` if invalid or missing.
- **Context**: Successfully validated tenant IDs are injected into the request context (`domain.WithTenantID`).
- **Usage**: Services and Repositories extract the tenant ID from the context to filter data.

## Soft Delete Policy

FinTrack Core implements a comprehensive soft delete strategy to maintain data auditability and prevent accidental data loss.

- **Traceable Deletion**: Critical entities (`Accounts`, `Transactions`, `Attachments`) support full traceability by storing both the timestamp (`deactivated_at`) and the ID of the user who performed the action (`deactivated_by`).
- **Standard Soft Delete**: Operational entities (`Users`, `Tenants`, `Tags`, `Categories`) use a standard `deactivated_at` timestamp.
- **Join Table Policy**: Many-to-many associations like `Users <-> Tenants` are soft-deleted via timestamp, while lightweight associations like `Transactions <-> Tags` are hard-deleted if they lack specific audit requirements in the schema.
- **Repository Pattern**: All "Read" operations (`Get`, `List`) automatically filter out soft-deleted records (`WHERE deactivated_at IS NULL`). "Delete" operations set the `deactivated_at` timestamp instead of removing the row.

## Data Integrity

The domain layer enforces business rules and data integrity through explicit `IsValid()` methods on all entities. This ensures that only valid data (e.g., non-negative balances, required fields, correct types) reaches the persistence layer.

### Environment Variables

Create a `.env` file in your root directory:

```env
DATABASE_URL=postgres://user:password@localhost:5432/fintrack?sslmode=disable
DB_HOST=localhost
DB_PORT=5432
DB_NAME=fintrack
DB_USER=postgres
DB_PASSWORD=postgres
SUPABASE_PROJECT_REF=your_supabase_project_ref
SUPABASE_ANON_KEY=your_supabase_anon_key
```

## Testing

The project uses `pgxmock` for unit testing the repository layer without requiring a live database.

```bash
go test ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.