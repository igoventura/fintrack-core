package dto

import (
	"time"

	"github.com/igoventura/fintrack-api/domain"
)

// CreateTransactionRequest represents the payload for creating a transaction.
type CreateTransactionRequest struct {
	FromAccountID   string                 `json:"from_account_id" binding:"required,uuid"`
	ToAccountID     *string                `json:"to_account_id,omitempty" binding:"omitempty,uuid"`
	Amount          float64                `json:"amount" binding:"required,gt=0"`
	AccrualMonth    string                 `json:"accrual_month" binding:"required,len=6"` // YYYYMM
	TransactionType domain.TransactionType `json:"transaction_type" binding:"required,oneof=credit debit transfer payment"`
	CategoryID      string                 `json:"category_id" binding:"required,uuid"`
	Comments        *string                `json:"comments,omitempty"`
	DueDate         time.Time              `json:"due_date" binding:"required"`
	PaymentDate     *time.Time             `json:"payment_date,omitempty"`
	TagIDs          []string               `json:"tag_ids,omitempty" binding:"omitempty,dive,uuid"`
}

// TransactionResponse represents the API response for a transaction.
type TransactionResponse struct {
	ID                  string                 `json:"id"`
	ParentTransactionID *string                `json:"parent_transaction_id,omitempty"`
	TenantID            string                 `json:"tenant_id"`
	FromAccountID       string                 `json:"from_account_id"`
	ToAccountID         *string                `json:"to_account_id,omitempty"`
	Currency            string                 `json:"currency"`
	Amount              float64                `json:"amount"`
	AccrualMonth        string                 `json:"accrual_month"`
	TransactionType     domain.TransactionType `json:"transaction_type"`
	CategoryID          string                 `json:"category_id"`
	Comments            *string                `json:"comments,omitempty"`
	DueDate             time.Time              `json:"due_date"`
	PaymentDate         *time.Time             `json:"payment_date,omitempty"`
	CreatedAt           time.Time              `json:"created_at"`
	CreatedBy           string                 `json:"created_by"`
	UpdatedAt           time.Time              `json:"updated_at"`
	UpdatedBy           string                 `json:"updated_by"`
	DeactivatedAt       *time.Time             `json:"deactivated_at,omitempty"`
	DeactivatedBy       *string                `json:"deactivated_by,omitempty"`
}

// ToDomain maps CreateTransactionRequest to domain.Transaction.
func (req *CreateTransactionRequest) ToDomain() *domain.Transaction {
	return &domain.Transaction{
		FromAccountID:   req.FromAccountID,
		ToAccountID:     req.ToAccountID,
		Amount:          req.Amount,
		AccrualMonth:    req.AccrualMonth,
		TransactionType: req.TransactionType,
		CategoryID:      req.CategoryID,
		Comments:        req.Comments,
		DueDate:         req.DueDate,
		PaymentDate:     req.PaymentDate,
	}
}

// FromTransactionDomain maps domain.Transaction to TransactionResponse.
func FromTransactionDomain(t *domain.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:                  t.ID,
		ParentTransactionID: t.ParentTransactionID,
		TenantID:            t.TenantID,
		FromAccountID:       t.FromAccountID,
		ToAccountID:         t.ToAccountID,
		Currency:            t.Currency,
		Amount:              t.Amount,
		AccrualMonth:        t.AccrualMonth,
		TransactionType:     t.TransactionType,
		CategoryID:          t.CategoryID,
		Comments:            t.Comments,
		DueDate:             t.DueDate,
		PaymentDate:         t.PaymentDate,
		CreatedAt:           t.CreatedAt,
		CreatedBy:           t.CreatedBy,
		UpdatedAt:           t.UpdatedAt,
		UpdatedBy:           t.UpdatedBy,
		DeactivatedAt:       t.DeactivatedAt,
		DeactivatedBy:       t.DeactivatedBy,
	}
}

// TransactionFilterRequest defines query parameters for listing transactions.
type TransactionFilterRequest struct {
	AccrualMonth    string                 `form:"accrual_month" binding:"omitempty,len=6"`
	AccountID       string                 `form:"account_id" binding:"omitempty,uuid"`
	TransactionType domain.TransactionType `form:"transaction_type" binding:"omitempty,oneof=credit debit transfer payment"`
}

// ToDomain maps TransactionFilterRequest to domain.TransactionFilter.
func (f *TransactionFilterRequest) ToDomain() domain.TransactionFilter {
	return domain.TransactionFilter{
		AccrualMonth:    f.AccrualMonth,
		AccountID:       f.AccountID,
		TransactionType: f.TransactionType,
	}
}
