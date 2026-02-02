# FinTrack Core

FinTrack Core is a Go library providing the foundational domain models and persistence logic for the FinTrack ecosystem. It is designed with modularity and testability in mind, following a domain-driven approach.

## Structure

```
fintrack-core/
├── domain/                # Business entities and repository interfaces
│   └── account.go         # Account model and AccountRepository interface
├── migrations/            # SQL migrations for tern
│   ├── tern.conf          # Tern configuration
│   └── *.sql              # Migration files
├── internal/
│   └── db/
│       └── postgres/      # PostgreSQL implementation using pgx
│           ├── db.go            # Connection pooling (pgxpool)
│           ├── account_repo.go  # PG implementation of AccountRepository
│           └── account_repo_test.go
├── go.mod                 # Go module definition
└── go.sum                 # Go module checksums
```

## Migrations

This project uses [tern](https://github.com/jackc/tern) for database migrations.

### Configuration

Edit `migrations/tern.conf` or use environment variables (`DB_HOST`, `DB_NAME`, etc.).

### Commands

```bash
# Run migrations
tern migrate -m migrations/

# Rollback last migration
tern rollback -m migrations/
```

## Getting Started

### Prerequisites

- Go 1.25.6 or higher
- PostgreSQL (for local development/runtime)

### Installation

```bash
go get github.com/igoventura/fintrack-core
```

### Usage

Initialize the PostgreSQL connection and the repository:

```go
package main

import (
	"context"
	"log"

	"github.com/igoventura/fintrack-core/internal/db/postgres"
)

func main() {
	ctx := context.Background()
	
	// Initialize Connection Pool
	connStr := "postgres://user:password@localhost:5432/fintrack?sslmode=disable"
	db, err := postgres.NewDB(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize Repository
	repo := postgres.NewAccountRepository(db.Pool)

	// Use the repository...
	accounts, err := repo.List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found %d accounts", len(accounts))
}
```

## Testing

The project uses `pgxmock` for unit testing the repository layer without requiring a live database.

```bash
go test ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.