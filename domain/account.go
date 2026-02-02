package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrAccountNotFound = errors.New("account not found")
)

// Account represents a financial account.
type Account struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AccountRepository defines the interface for account persistence.
type AccountRepository interface {
	GetByID(ctx context.Context, id string) (*Account, error)
	List(ctx context.Context) ([]Account, error)
	Create(ctx context.Context, acc *Account) error
	Update(ctx context.Context, acc *Account) error
	Delete(ctx context.Context, id string) error
}
