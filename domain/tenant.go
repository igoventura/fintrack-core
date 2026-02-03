package domain

import (
	"context"
	"errors"
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

type contextKey string

const tenantIDKey contextKey = "tenantID"

// WithTenantID returns a new context with the given tenant ID.
func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

// GetTenantID retrieves the tenant ID from the context.
func GetTenantID(ctx context.Context) string {
	val, _ := ctx.Value(tenantIDKey).(string)
	return val
}

func (t *Tenant) IsValid() (bool, map[string]error) {
	err := make(map[string]error)
	if t.Name == "" {
		err["name"] = errors.New("name is required")
	}
	if len(err) == 0 {
		return true, nil
	}
	return false, err
}
