package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/igoventura/fintrack-api/domain"
)

type TransactionService struct {
	repo         domain.TransactionRepository
	accountRepo  domain.AccountRepository
	categoryRepo domain.CategoryRepository
	tagRepo      domain.TagRepository
}

func NewTransactionService(
	repo domain.TransactionRepository,
	accountRepo domain.AccountRepository,
	categoryRepo domain.CategoryRepository,
	tagRepo domain.TagRepository,
) *TransactionService {
	return &TransactionService{
		repo:         repo,
		accountRepo:  accountRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
	}
}

func (s *TransactionService) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	tenantID := domain.GetTenantID(ctx)
	return s.repo.GetByID(ctx, tenantID, id)
}

func (s *TransactionService) List(ctx context.Context, filter domain.TransactionFilter) ([]domain.Transaction, error) {
	tenantID := domain.GetTenantID(ctx)
	return s.repo.List(ctx, tenantID, filter)
}

func (s *TransactionService) Create(ctx context.Context, t *domain.Transaction, tagIDs []string, installments int, isRecurring bool) error {
	tenantID := domain.GetTenantID(ctx)
	t.TenantID = tenantID

	userID := domain.GetUserID(ctx)
	if userID == "" {
		return errors.New("user ID is required")
	}
	t.CreatedBy = userID
	t.UpdatedBy = userID

	// 1. Fetch FromAccount (Needed for Currency and PaymentDate logic)
	fromAccount, err := s.accountRepo.GetByID(ctx, t.FromAccountID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to fetch from_account: %w", err)
	}
	if fromAccount.TenantID != tenantID {
		return errors.New("from_account does not belong to this tenant")
	}

	// 2. Field Defaults
	// Currency: Default to FromAccount currency if not set
	if t.Currency == "" {
		t.Currency = fromAccount.Currency
	}

	// AccrualMonth: Default to DueDate YYYYMM if not set
	if t.AccrualMonth == "" {
		t.AccrualMonth = t.DueDate.Format("200601")
	}

	// PaymentDate: Default to DueDate if Credit Card and type is Debit/Credit
	if t.PaymentDate == nil {
		if fromAccount.Type == domain.AccountTypeCreditCard &&
			(t.TransactionType == domain.TransactionTypeCredit || t.TransactionType == domain.TransactionTypeDebit) {
			t.PaymentDate = &t.DueDate
		}
	}

	// Validate basic fields (Now that defaults are set)
	if valid, errs := t.IsValid(); !valid {
		var errMsg string
		for field, err := range errs {
			errMsg += fmt.Sprintf("%s: %s; ", field, err.Error())
		}
		return errors.New("validation failed: " + errMsg)
	}

	// 3. ToAccount Validation (if applicable)
	if t.ToAccountID != nil && *t.ToAccountID != "" {
		toAccount, err := s.accountRepo.GetByID(ctx, *t.ToAccountID, tenantID)
		if err != nil {
			return fmt.Errorf("failed to fetch to_account: %w", err)
		}
		if toAccount.TenantID != tenantID {
			return errors.New("to_account does not belong to this tenant")
		}
	}

	// 4. Category Validation
	if _, err := s.categoryRepo.GetByID(ctx, t.CategoryID, tenantID); err != nil {
		return fmt.Errorf("invalid category: %w", err)
	}

	// 5. Tags Validation
	if len(tagIDs) > 0 {
		valid, err := s.tagRepo.ValidateTags(ctx, tenantID, tagIDs)
		if err != nil {
			return fmt.Errorf("failed to validate tags: %w", err)
		}
		if !valid {
			return errors.New("one or more tags do not belong to this tenant")
		}
	}

	// 6. Installments Logic
	numInstallments := installments
	if numInstallments < 1 {
		numInstallments = 1
	}

	if numInstallments > 1 {
		calcInstallments, err := domain.CalculateInstallments(t.Amount, numInstallments, t.DueDate, isRecurring)
		if err != nil {
			return fmt.Errorf("failed to calculate installments: %w", err)
		}

		// Prepare Parent (1st Installment)
		t.Amount = calcInstallments[0].Amount
		t.DueDate = calcInstallments[0].DueDate

		// Recalculate AccrualMonth based on the new DueDate
		t.AccrualMonth = t.DueDate.Format("200601")

		var comments string
		if t.Comments != nil {
			comments = *t.Comments
		}
		newComments := fmt.Sprintf("[Installment %d/%d] %s", calcInstallments[0].Number, numInstallments, comments)
		t.Comments = &newComments

		// Prepare Children
		children := make([]domain.Transaction, 0, numInstallments-1)
		for i := 1; i < numInstallments; i++ {
			inst := calcInstallments[i]
			child := *t   // Copy Base
			child.ID = "" // Clear ID (will be generated)
			child.Amount = inst.Amount
			child.DueDate = inst.DueDate
			child.AccrualMonth = child.DueDate.Format("200601")

			// PaymentDate logic for children? same rules apply.
			if fromAccount.Type == domain.AccountTypeCreditCard &&
				(t.TransactionType == domain.TransactionTypeCredit || t.TransactionType == domain.TransactionTypeDebit) {
				child.PaymentDate = &child.DueDate
			} else {
				// Future installments are generally unpaid, so clear PaymentDate unless explicitly handled above.
				if t.PaymentDate != nil {
					child.PaymentDate = nil
				}
			}
			// Apply rule again just in case
			if fromAccount.Type == domain.AccountTypeCreditCard &&
				(t.TransactionType == domain.TransactionTypeCredit || t.TransactionType == domain.TransactionTypeDebit) {
				child.PaymentDate = &child.DueDate
			}

			childComments := fmt.Sprintf("[Installment %d/%d] %s", inst.Number, numInstallments, comments)
			child.Comments = &childComments

			// Tenant/User are copied
			child.ParentTransactionID = nil // Will be linked in Repo

			children = append(children, child)
		}

		if err := s.repo.CreateWithInstallments(ctx, t, children, tagIDs); err != nil {
			return fmt.Errorf("failed to create transaction with installments: %w", err)
		}

	} else {
		// Single Transaction
		if err := s.repo.Create(ctx, t); err != nil {
			return err
		}

		// Tags Association
		if len(tagIDs) > 0 {
			if err := s.repo.AddTagsToTransaction(ctx, t.ID, tagIDs); err != nil {
				return fmt.Errorf("transaction created but failed to link tags: %w", err)
			}
		}
	}

	return nil
}

func (s *TransactionService) Update(ctx context.Context, t *domain.Transaction, tagIDs []string) error {
	tenantID := domain.GetTenantID(ctx)
	t.TenantID = tenantID // Ensure we don't overwrite with wrong tenant

	userID := domain.GetUserID(ctx)
	if userID == "" {
		return errors.New("user ID is required")
	}
	t.UpdatedBy = userID

	// Validate tags if provided
	if len(tagIDs) > 0 {
		valid, err := s.tagRepo.ValidateTags(ctx, tenantID, tagIDs)
		if err != nil {
			return fmt.Errorf("failed to validate tags: %w", err)
		}
		if !valid {
			return errors.New("one or more tags do not belong to this tenant")
		}
	}

	if err := s.repo.Update(ctx, t); err != nil {
		return err
	}

	// Update tags (Replace strategy)
	if tagIDs != nil {
		if err := s.repo.ReplaceTags(ctx, t.ID, tagIDs); err != nil {
			return fmt.Errorf("transaction updated but failed to update tags: %w", err)
		}
	}

	return nil
}

func (s *TransactionService) Delete(ctx context.Context, id string) error {
	tenantID := domain.GetTenantID(ctx)
	userID := domain.GetUserID(ctx)
	if userID == "" {
		return errors.New("user ID is required")
	}
	return s.repo.Delete(ctx, tenantID, id, userID)
}
