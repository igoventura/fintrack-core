package domain

import (
	"context"
	"time"
)

// Tenant represents a tenant in the system.
type Tenant struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivated_at,omitempty"`
}

// TenantRepository defines the interface for tenant persistence.
type TenantRepository interface {
	GetByID(ctx context.Context, id string) (*Tenant, error)
	Create(ctx context.Context, tenant *Tenant) error
	Update(ctx context.Context, tenant *Tenant) error
	Delete(ctx context.Context, id string) error
	ListByUserID(ctx context.Context, userID string) ([]Tenant, error)
}
