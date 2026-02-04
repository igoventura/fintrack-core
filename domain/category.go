package domain

import (
	"context"
	"errors"
	"time"
)

// Category represents a classification for transactions.
type Category struct {
	ID               string     `json:"id"`
	ParentCategoryID *string    `json:"parent_category_id,omitempty"`
	TenantID         string     `json:"tenant_id"`
	Name             string     `json:"name"`
	DeactivatedAt    *time.Time `json:"deactivated_at,omitempty"`
	Color            string     `json:"color"`
	Icon             string     `json:"icon"`
	CreatedAt        time.Time  `json:"created_at"`
	CreatedBy        string     `json:"created_by"`
	UpdatedAt        time.Time  `json:"updated_at"`
	UpdatedBy        string     `json:"updated_by"`
	DeactivatedBy    *string    `json:"deactivated_by,omitempty"`
}

// CategoryRepository defines the interface for category persistence.
type CategoryRepository interface {
	GetByID(ctx context.Context, id, tenantID string) (*Category, error)
	List(ctx context.Context, tenantID string) ([]Category, error)
	Create(ctx context.Context, cat *Category) error
	Update(ctx context.Context, cat *Category) error
	Delete(ctx context.Context, id, tenantID, userID string) error
}

func (c *Category) IsValid() (bool, map[string]error) {
	err := make(map[string]error)
	if c.Name == "" {
		err["name"] = errors.New("name is required")
	}
	if c.TenantID == "" {
		err["tenant_id"] = errors.New("tenant_id is required")
	}
	if len(err) == 0 {
		return true, nil
	}
	return false, err
}
