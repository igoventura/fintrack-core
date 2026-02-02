package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrAccountNotFound = errors.New("account not found")
)

// AccountType represents the type of a financial account.
type AccountType string

const (
	AccountTypeBank       AccountType = "bank"
	AccountTypeCash       AccountType = "cash"
	AccountTypeCreditCard AccountType = "credit_card"
	AccountTypeInvestment AccountType = "investment"
	AccountTypeOther      AccountType = "other"
)

// CreditCardBrand represents the brand of a credit card.
type CreditCardBrand string

const (
	BrandVisa       CreditCardBrand = "visa"
	BrandMastercard CreditCardBrand = "mastercard"
	BrandAmex       CreditCardBrand = "amex"
	BrandDiscover   CreditCardBrand = "discover"
	BrandJCB        CreditCardBrand = "jcb"
	BrandUnionpay   CreditCardBrand = "unionpay"
	BrandDinersClub CreditCardBrand = "diners_club"
	BrandMaestro    CreditCardBrand = "maestro"
	BrandUnknown    CreditCardBrand = "unknown"
)

// Account represents a financial account.
type Account struct {
	ID             string      `json:"id"`
	TenantID       string      `json:"tenant_id"`
	Name           string      `json:"name"`
	InitialBalance float64     `json:"initial_balance"`
	Color          string      `json:"color"`
	Currency       string      `json:"currency"`
	Icon           string      `json:"icon"`
	Type           AccountType `json:"type"`
	CreatedAt      time.Time   `json:"created_at"`
	CreatedBy      string      `json:"created_by"`
	UpdatedAt      time.Time   `json:"updated_at"`
	UpdatedBy      string      `json:"updated_by"`
	DeactivatedAt  *time.Time  `json:"deactivated_at,omitempty"`
	DeactivatedBy  *string     `json:"deactivated_by,omitempty"`
}

// CreditCardInfo represents detailed information for a credit card account.
type CreditCardInfo struct {
	ID            string          `json:"id"`
	AccountID     string          `json:"account_id"`
	LastFour      string          `json:"last_four"`
	Name          string          `json:"name"`
	Brand         CreditCardBrand `json:"brand"`
	ClosingDate   time.Time       `json:"closing_date"`
	DueDate       time.Time       `json:"due_date"`
	CreatedAt     time.Time       `json:"created_at"`
	CreatedBy     string          `json:"created_by"`
	UpdatedAt     time.Time       `json:"updated_at"`
	UpdatedBy     string          `json:"updated_by"`
	DeactivatedAt *time.Time      `json:"deactivated_at,omitempty"`
	DeactivatedBy *string         `json:"deactivated_by,omitempty"`
}

// AccountRepository defines the interface for account persistence.
type AccountRepository interface {
	GetByID(ctx context.Context, id string) (*Account, error)
	List(ctx context.Context, tenantID string) ([]Account, error)
	Create(ctx context.Context, acc *Account) error
	Update(ctx context.Context, acc *Account) error
	Delete(ctx context.Context, id string, userID string) error

	GetCreditCardInfo(ctx context.Context, accountID string) (*CreditCardInfo, error)
	UpsertCreditCardInfo(ctx context.Context, info *CreditCardInfo) error
}
