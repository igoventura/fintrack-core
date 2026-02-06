package dto

import "time"

type CreateCategoryRequest struct {
	Name             string  `json:"name" binding:"required"`
	Type             string  `json:"type" binding:"required,oneof=income expense transfer"`
	ParentCategoryID *string `json:"parent_category_id,omitempty"`
	Color            string  `json:"color"`
	Icon             string  `json:"icon"`
}

type UpdateCategoryRequest struct {
	Name             string  `json:"name" binding:"required"`
	ParentCategoryID *string `json:"parent_category_id,omitempty"`
	Color            string  `json:"color"`
	Icon             string  `json:"icon"`
}

type CategoryResponse struct {
	ID               string    `json:"id"`
	ParentCategoryID *string   `json:"parent_category_id,omitempty"`
	TenantID         string    `json:"tenant_id"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`
	Color            string    `json:"color"`
	Icon             string    `json:"icon"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        string    `json:"created_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        string    `json:"updated_by"`
}
