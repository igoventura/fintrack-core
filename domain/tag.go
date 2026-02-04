package domain

import (
	"context"
	"errors"
	"time"
)

// Tag represents a label for transactions.
type Tag struct {
	ID            string     `json:"id"`
	TenantID      string     `json:"tenant_id"`
	Name          string     `json:"name"`
	CreatedAt     time.Time  `json:"created_at"`
	CreatedBy     string     `json:"created_by"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UpdatedBy     string     `json:"updated_by"`
	DeactivatedAt *time.Time `json:"deactivated_at,omitempty"`
	DeactivatedBy *string    `json:"deactivated_by,omitempty"`
}

// TagRepository defines the interface for tag persistence.
type TagRepository interface {
	GetByID(ctx context.Context, id, tenantID string) (*Tag, error)
	List(ctx context.Context, tenantID string) ([]Tag, error)
	Create(ctx context.Context, tag *Tag) error
	Update(ctx context.Context, tag *Tag) error
	Delete(ctx context.Context, id, tenantID, userID string) error
	ValidateTags(ctx context.Context, tenantID string, tagIDs []string) (bool, error)
}

func (t *Tag) IsValid() (bool, map[string]error) {
	err := make(map[string]error)
	if t.Name == "" {
		err["name"] = errors.New("name is required")
	}
	if t.TenantID == "" {
		err["tenant_id"] = errors.New("tenant_id is required")
	}
	if len(err) == 0 {
		return true, nil
	}
	return false, err
}
