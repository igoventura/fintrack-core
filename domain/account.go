package domain

import (
	"context"
	"errors"
	"slices"
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
	GetByID(ctx context.Context, id, tenantID string) (*Account, error)
	List(ctx context.Context, tenantID string) ([]Account, error)
	Create(ctx context.Context, acc *Account) error
	Update(ctx context.Context, acc *Account) error
	Delete(ctx context.Context, id, tenantID, userID string) error

	GetCreditCardInfo(ctx context.Context, accountID string) (*CreditCardInfo, error)
	UpsertCreditCardInfo(ctx context.Context, info *CreditCardInfo) error
}

func (a *Account) IsValid() (bool, map[string]error) {
	err := make(map[string]error)
	if a.Name == "" {
		err["name"] = errors.New("name is required")
	}
	if a.TenantID == "" {
		err["tenant_id"] = errors.New("tenant_id is required")
	}
	if a.InitialBalance < 0 {
		err["initial_balance"] = errors.New("initial_balance must be non-negative")
	}
	if a.Currency == "" {
		err["currency"] = errors.New("currency is required")
	}
	if a.Color == "" {
		err["color"] = errors.New("color is required")
	} else if len(a.Color) > 128 {
		err["color"] = errors.New("color must not exceed 128 characters")
	}
	if a.Type == "" {
		err["type"] = errors.New("type is required")
	} else {
		validTypes := []AccountType{AccountTypeBank, AccountTypeCash, AccountTypeCreditCard, AccountTypeInvestment, AccountTypeOther}
		if !slices.Contains(validTypes, a.Type) {
			err["type"] = errors.New("invalid account type")
		}
	}
	if len(err) == 0 {
		return true, nil
	}
	return false, err
}

func (cci *CreditCardInfo) IsValid() (bool, map[string]error) {
	err := make(map[string]error)
	if cci.AccountID == "" {
		err["account_id"] = errors.New("account_id is required")
	}
	if cci.LastFour == "" {
		err["last_four"] = errors.New("last_four is required")
	}
	if cci.Name == "" {
		err["name"] = errors.New("name is required")
	}
	if cci.Brand == "" {
		err["brand"] = errors.New("brand is required")
	} else {
		validBrands := []CreditCardBrand{BrandVisa, BrandMastercard, BrandAmex, BrandDiscover, BrandJCB, BrandUnionpay, BrandDinersClub, BrandMaestro, BrandUnknown}
		if !slices.Contains(validBrands, cci.Brand) {
			err["brand"] = errors.New("invalid brand")
		}
	}
	if cci.ClosingDate.IsZero() {
		err["closing_date"] = errors.New("closing_date is required")
	}
	if cci.DueDate.IsZero() {
		err["due_date"] = errors.New("due_date is required")
	}
	if len(err) == 0 {
		return true, nil
	}
	return false, err
}
