package domain

import (
	"context"
	"errors"
	"slices"
	"time"
)

// TransactionType represents the type of a transaction.
type TransactionType string

const (
	TransactionTypeCredit   TransactionType = "credit"
	TransactionTypeDebit    TransactionType = "debit"
	TransactionTypeTransfer TransactionType = "transfer"
	TransactionTypePayment  TransactionType = "payment"
)

// Transaction represents a financial movement.
type Transaction struct {
	ID                           string          `json:"id"`
	PreviousSiblingTransactionID *string         `json:"previous_sibling_transaction_id,omitempty"`
	NextSiblingTransactionID     *string         `json:"next_sibling_transaction_id,omitempty"`
	TenantID                     string          `json:"tenant_id"`
	FromAccountID                string          `json:"from_account_id"`
	ToAccountID                  *string         `json:"to_account_id,omitempty"`
	Amount                       float64         `json:"amount"`
	AccrualMonth                 string          `json:"accrual_month"` // YYYYMM
	TransactionType              TransactionType `json:"transaction_type"`
	CategoryID                   string          `json:"category_id"`
	Comments                     *string         `json:"comments,omitempty"`
	DueDate                      time.Time       `json:"due_date"`
	PaymentDate                  *time.Time      `json:"payment_date,omitempty"`
	CreatedAt                    time.Time       `json:"created_at"`
	CreatedBy                    string          `json:"created_by"`
	UpdatedAt                    time.Time       `json:"updated_at"`
	UpdatedBy                    string          `json:"updated_by"`
	DeactivatedAt                *time.Time      `json:"deactivated_at,omitempty"`
	DeactivatedBy                *string         `json:"deactivated_by,omitempty"`
}

// TransactionTag represents the association between a transaction and a tag.
type TransactionTag struct {
	TransactionID string `json:"transaction_id"`
	TagID         string `json:"tag_id"`
}

// TransactionAttachment represents a file attached to a transaction.
type TransactionAttachment struct {
	ID            string     `json:"id"`
	TransactionID string     `json:"transaction_id"`
	Name          string     `json:"name"`
	Path          string     `json:"path"`
	CreatedAt     time.Time  `json:"created_at"`
	CreatedBy     string     `json:"created_by"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UpdatedBy     string     `json:"updated_by"`
	DeactivatedAt *time.Time `json:"deactivated_at,omitempty"`
	DeactivatedBy *string    `json:"deactivated_by,omitempty"`
}

// TransactionRepository defines the interface for transaction persistence.
type TransactionRepository interface {
	GetByID(ctx context.Context, id string) (*Transaction, error)
	List(ctx context.Context, tenantID string) ([]Transaction, error)
	Create(ctx context.Context, tx *Transaction) error
	Update(ctx context.Context, tx *Transaction) error
	Delete(ctx context.Context, id string, userID string) error

	// Tag associations
	AddTagToTransaction(ctx context.Context, transactionID, tagID string) error
	RemoveTagFromTransaction(ctx context.Context, transactionID, tagID string) error
	ListTransactionTags(ctx context.Context, transactionID string) ([]Tag, error)

	// Attachment associations
	AddAttachment(ctx context.Context, attachment *TransactionAttachment) error
	RemoveAttachment(ctx context.Context, id string, userID string) error
	ListAttachments(ctx context.Context, transactionID string) ([]TransactionAttachment, error)
}

func (t *Transaction) IsValid() (bool, map[string]error) {
	err := make(map[string]error)
	if t.TenantID == "" {
		err["tenant_id"] = errors.New("tenant_id is required")
	}
	if t.FromAccountID == "" {
		err["from_account_id"] = errors.New("from_account_id is required")
	}
	if t.Amount <= 0 {
		err["amount"] = errors.New("amount must be greater than 0")
	}
	if t.TransactionType == "" {
		err["transaction_type"] = errors.New("transaction_type is required")
	} else {
		validTypes := []TransactionType{TransactionTypeCredit, TransactionTypeDebit, TransactionTypeTransfer, TransactionTypePayment}
		if !slices.Contains(validTypes, t.TransactionType) {
			err["transaction_type"] = errors.New("invalid transaction type")
		} else if t.TransactionType == TransactionTypeTransfer {
			if t.ToAccountID == nil || *t.ToAccountID == "" {
				err["to_account_id"] = errors.New("to_account_id is required for transfers")
			}
			if t.FromAccountID == *t.ToAccountID {
				err["to_account_id"] = errors.New("to_account_id must be different from from_account_id")
			}
		}
	}
	if t.CategoryID == "" {
		err["category_id"] = errors.New("category_id is required")
	}
	if len(t.AccrualMonth) != 6 {
		// Simple length check for YYYYMM. Could be more robust with regex but good baseline.
		err["accrual_month"] = errors.New("accrual_month must be in YYYYMM format")
	}
	if t.DueDate.IsZero() {
		err["due_date"] = errors.New("due_date is required")
	}
	if len(err) == 0 {
		return true, nil
	}
	return false, err
}
