package domain

import (
	"context"
	"errors"
	"regexp"
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

const userIdKey contextKey = "userID"

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIdKey, userID)
}

func GetUserID(ctx context.Context) string {
	return ctx.Value(userIdKey).(string)
}

func (u *User) IsValid() (bool, map[string]error) {
	err := make(map[string]error)
	if u.Name == "" {
		err["name"] = errors.New("name is required")
	}
	if u.Email == "" {
		err["email"] = errors.New("email is required")
	} else {
		// Basic email format validation
		const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		matched, _ := regexp.MatchString(emailRegexPattern, u.Email)
		if !matched {
			err["email"] = errors.New("invalid email format")
		}
	}
	if u.SupabaseID == "" {
		err["supabase_id"] = errors.New("supabase_id is required")
	}
	if len(err) == 0 {
		return true, nil
	}
	return false, err
}
