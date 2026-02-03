package domain

import (
	"context"
	"time"
)

// User represents a user in the system.
type User struct {
	ID            string     `json:"id"`
	SupabaseID    string     `json:"supabase_id"`
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivated_at,omitempty"`
}

// UserTenant represents the association between a user and a tenant.
type UserTenant struct {
	UserID        string     `json:"user_id"`
	TenantID      string     `json:"tenant_id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivated_at,omitempty"`
}

// UserRepository defines the interface for user persistence.
type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetBySupabaseID(ctx context.Context, supabaseID string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error

	// Tenant associations
	AddUserToTenant(ctx context.Context, userID, tenantID string) error
	RemoveUserFromTenant(ctx context.Context, userID, tenantID string) error
	ListUserTenants(ctx context.Context, userID string) ([]UserTenant, error)
}
