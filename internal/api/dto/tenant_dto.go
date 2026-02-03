package dto

import "time"

// CreateTenantRequest represents the payload for creating a new tenant.
type CreateTenantRequest struct {
	Name string `json:"name" binding:"required"`
}

// TenantResponse represents a tenant in API responses.
type TenantResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
