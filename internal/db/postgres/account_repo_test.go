package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/igoventura/fintrack-core/domain"
	"github.com/pashagolub/pgxmock/v3"
)

func TestAccountRepository(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := NewAccountRepository(mock)
	ctx := context.Background()

	t.Run("GetByID", func(t *testing.T) {
		now := time.Now()
		rows := pgxmock.NewRows([]string{"id", "name", "balance", "created_at", "updated_at"}).
			AddRow("acc-1", "Savings", 1000.0, now, now)

		mock.ExpectQuery("SELECT id, name, balance, created_at, updated_at FROM accounts WHERE id = \\$1").
			WithArgs("acc-1").
			WillReturnRows(rows)

		acc, err := repo.GetByID(ctx, "acc-1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if acc.Name != "Savings" {
			t.Errorf("expected Savings, got %s", acc.Name)
		}
	})

	t.Run("Create", func(t *testing.T) {
		acc := &domain.Account{
			ID:      "acc-2",
			Name:    "Checking",
			Balance: 500.0,
		}

		mock.ExpectExec("INSERT INTO accounts").
			WithArgs("acc-2", "Checking", 500.0, pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := repo.Create(ctx, acc)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
