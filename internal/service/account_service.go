package service

import (
	"context"
	"fmt"

	"github.com/igoventura/fintrack-core/domain"
)

type AccountService struct {
	repo domain.AccountRepository
}

func NewAccountService(repo domain.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	tenantID := domain.GetTenantID(ctx)
	acc, err := s.repo.GetByID(ctx, id, tenantID)
	if err != nil {
		return nil, fmt.Errorf("service failed to get account: %w", err)
	}
	return acc, nil
}

func (s *AccountService) ListAccounts(ctx context.Context) ([]domain.Account, error) {
	tenantID := domain.GetTenantID(ctx)
	accounts, err := s.repo.List(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("service failed to list accounts: %w", err)
	}
	return accounts, nil
}

func (s *AccountService) CreateAccount(ctx context.Context, acc *domain.Account) error {
	tenantID := domain.GetTenantID(ctx)
	acc.TenantID = tenantID

	// Business validation logic could go here
	if acc.Name == "" {
		return fmt.Errorf("account name is required")
	}

	if err := s.repo.Create(ctx, acc); err != nil {
		return fmt.Errorf("service failed to create account: %w", err)
	}
	return nil
}

func (s *AccountService) UpdateAccount(ctx context.Context, acc *domain.Account) error {
	tenantID := domain.GetTenantID(ctx)
	acc.TenantID = tenantID

	if err := s.repo.Update(ctx, acc); err != nil {
		return fmt.Errorf("service failed to update account: %w", err)
	}
	return nil
}

func (s *AccountService) DeleteAccount(ctx context.Context, id string, userID string) error {
	tenantID := domain.GetTenantID(ctx)

	if err := s.repo.Delete(ctx, id, tenantID, userID); err != nil {
		return fmt.Errorf("service failed to delete account: %w", err)
	}
	return nil
}
