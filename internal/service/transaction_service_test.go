package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/igoventura/fintrack-api/domain"
)

// Mocks
type mockAccountRepo struct {
	domain.AccountRepository
	GetByIDFn func(ctx context.Context, id, tenantID string) (*domain.Account, error)
}

func (m *mockAccountRepo) GetByID(ctx context.Context, id, tenantID string) (*domain.Account, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id, tenantID)
	}
	return nil, errors.New("not implemented")
}

type mockRepo struct {
	domain.TransactionRepository
	CreateFn                 func(ctx context.Context, tx *domain.Transaction) error
	CreateWithInstallmentsFn func(ctx context.Context, parent *domain.Transaction, children []domain.Transaction, tagIDs []string) error
}

func (m *mockRepo) Create(ctx context.Context, tx *domain.Transaction) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, tx)
	}
	return nil
}

func (m *mockRepo) CreateWithInstallments(ctx context.Context, parent *domain.Transaction, children []domain.Transaction, tagIDs []string) error {
	if m.CreateWithInstallmentsFn != nil {
		return m.CreateWithInstallmentsFn(ctx, parent, children, tagIDs)
	}
	return nil
}

type mockCategoryRepo struct {
	domain.CategoryRepository
	GetByIDFn func(ctx context.Context, id, tenantID string) (*domain.Category, error)
}

func (m *mockCategoryRepo) GetByID(ctx context.Context, id, tenantID string) (*domain.Category, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id, tenantID)
	}
	return &domain.Category{ID: id, TenantID: tenantID}, nil
}

type mockTagRepo struct {
	domain.TagRepository
	ValidateTagsFn func(ctx context.Context, tenantID string, tagIDs []string) (bool, error)
}

func (m *mockTagRepo) ValidateTags(ctx context.Context, tenantID string, tagIDs []string) (bool, error) {
	if m.ValidateTagsFn != nil {
		return m.ValidateTagsFn(ctx, tenantID, tagIDs)
	}
	return true, nil
}

func TestTransactionService_Create(t *testing.T) {
	ctx := context.Background()
	ctx = domain.WithTenantID(ctx, "tenant-1")
	ctx = domain.WithUserID(ctx, "user-1")

	toAccID := "acc-2"
	comments := "Test"

	tests := []struct {
		name         string
		transaction  *domain.Transaction
		installments int
		isRecurring  bool
		setupMocks   func(*mockRepo, *mockAccountRepo)
		expectError  bool
	}{
		{
			name: "Single Transaction - PaymentDate Logic (CC)",
			transaction: &domain.Transaction{
				FromAccountID:   "acc-cc",
				Amount:          100,
				TransactionType: domain.TransactionTypeCredit,
				Comments:        &comments,
				CategoryID:      "cat-1",
				DueDate:         time.Now(),
				TenantID:        "tenant-1",
			},
			installments: 1,
			setupMocks: func(r *mockRepo, ar *mockAccountRepo) {
				ar.GetByIDFn = func(ctx context.Context, id, tenantID string) (*domain.Account, error) {
					return &domain.Account{ID: id, TenantID: tenantID, Type: domain.AccountTypeCreditCard, Currency: "USD"}, nil
				}
				r.CreateFn = func(ctx context.Context, tx *domain.Transaction) error {
					if tx.PaymentDate == nil {
						t.Error("expected PaymentDate to be set for CC Credit")
					}
					if *tx.PaymentDate != tx.DueDate {
						t.Errorf("expected PaymentDate %v, got %v", tx.DueDate, *tx.PaymentDate)
					}
					return nil
				}
			},
			expectError: false,
		},
		{
			name: "Split Installments (3x)",
			transaction: &domain.Transaction{
				FromAccountID:   "acc-1",
				Amount:          100,
				TransactionType: domain.TransactionTypeDebit,
				CategoryID:      "cat-1",
				DueDate:         time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
				TenantID:        "tenant-1", // Will be overwritten by service but safe to set
				ToAccountID:     &toAccID,
			},
			installments: 3,
			isRecurring:  false,
			setupMocks: func(r *mockRepo, ar *mockAccountRepo) {
				ar.GetByIDFn = func(ctx context.Context, id, tenantID string) (*domain.Account, error) {
					if id == "acc-1" {
						return &domain.Account{ID: id, TenantID: tenantID, Type: domain.AccountTypeBank, Currency: "USD"}, nil
					}
					return &domain.Account{ID: id, TenantID: tenantID, Type: domain.AccountTypeBank, Currency: "USD"}, nil
				}
				r.CreateWithInstallmentsFn = func(ctx context.Context, parent *domain.Transaction, children []domain.Transaction, tagIDs []string) error {
					if len(children) != 2 {
						t.Errorf("expected 2 children, got %d", len(children))
					}
					// Parent: 33.34
					if parent.Amount != 33.34 {
						t.Errorf("parent amount mismatch: got %f, want 33.34", parent.Amount)
					}
					// Children: 33.33 each
					if children[0].Amount != 33.33 {
						t.Errorf("child 0 amount mismatch: got %f, want 33.33", children[0].Amount)
					}
					return nil
				}
			},
			expectError: false,
		},
		{
			name: "Recurring Installments (3x)",
			transaction: &domain.Transaction{
				FromAccountID:   "acc-1",
				Amount:          100,
				TransactionType: domain.TransactionTypeDebit,
				CategoryID:      "cat-1",
				DueDate:         time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			},
			installments: 3,
			isRecurring:  true,
			setupMocks: func(r *mockRepo, ar *mockAccountRepo) {
				ar.GetByIDFn = func(ctx context.Context, id, tenantID string) (*domain.Account, error) {
					return &domain.Account{ID: id, TenantID: tenantID, Type: domain.AccountTypeBank, Currency: "USD"}, nil
				}
				r.CreateWithInstallmentsFn = func(ctx context.Context, parent *domain.Transaction, children []domain.Transaction, tagIDs []string) error {
					// Parent: 100
					if parent.Amount != 100 {
						t.Errorf("parent amount mismatch: %f", parent.Amount)
					}
					// Children: 100
					if children[0].Amount != 100 {
						t.Errorf("child amount mismatch: %f", children[0].Amount)
					}
					return nil
				}
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{}
			accRepo := &mockAccountRepo{}
			catRepo := &mockCategoryRepo{}
			tagRepo := &mockTagRepo{}

			if tt.setupMocks != nil {
				tt.setupMocks(repo, accRepo)
			}

			s := NewTransactionService(repo, accRepo, catRepo, tagRepo)
			err := s.Create(ctx, tt.transaction, nil, tt.installments, tt.isRecurring)

			if (err != nil) != tt.expectError {
				t.Errorf("Create() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}
