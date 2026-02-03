package dto

import (
	"time"

	"github.com/igoventura/fintrack-core/domain"
)

type AccountResponse struct {
	ID             string             `json:"id"`
	TenantID       string             `json:"tenant_id"`
	Name           string             `json:"name"`
	InitialBalance float64            `json:"initial_balance"`
	Color          string             `json:"color"`
	Currency       string             `json:"currency"`
	Icon           string             `json:"icon"`
	Type           domain.AccountType `json:"type"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

func MapAccountToResponse(acc *domain.Account) *AccountResponse {
	return &AccountResponse{
		ID:             acc.ID,
		TenantID:       acc.TenantID,
		Name:           acc.Name,
		InitialBalance: acc.InitialBalance,
		Color:          acc.Color,
		Currency:       acc.Currency,
		Icon:           acc.Icon,
		Type:           acc.Type,
		CreatedAt:      acc.CreatedAt,
		UpdatedAt:      acc.UpdatedAt,
	}
}

type CreateAccountRequest struct {
	Name           string             `json:"name" validate:"required"`
	InitialBalance float64            `json:"initial_balance"`
	Color          string             `json:"color"`
	Currency       string             `json:"currency"`
	Icon           string             `json:"icon"`
	Type           domain.AccountType `json:"type" validate:"required"`
}

func (r *CreateAccountRequest) ToEntity(creatorID string, tenantID string) *domain.Account {
	now := time.Now()
	return &domain.Account{
		TenantID:       tenantID,
		Name:           r.Name,
		InitialBalance: r.InitialBalance,
		Color:          r.Color,
		Currency:       r.Currency,
		Icon:           r.Icon,
		Type:           r.Type,
		CreatedAt:      now,
		CreatedBy:      creatorID,
		UpdatedAt:      now,
		UpdatedBy:      creatorID,
	}
}
