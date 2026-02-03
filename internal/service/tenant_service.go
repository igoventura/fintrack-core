package service

import (
	"context"
	"fmt"

	"github.com/igoventura/fintrack-core/domain"
)

type TenantService struct {
	repo        domain.TenantRepository
	userService *UserService
}

func NewTenantService(repo domain.TenantRepository, userService *UserService) *TenantService {
	return &TenantService{
		repo:        repo,
		userService: userService,
	}
}

func (s *TenantService) CreateTenant(ctx context.Context, name, creatorID string) (*domain.Tenant, error) {
	if name == "" {
		return nil, fmt.Errorf("tenant name is required")
	}

	tenant := &domain.Tenant{
		Name: name,
	}

	// Create tenant
	if err := s.repo.Create(ctx, tenant); err != nil {
		return nil, fmt.Errorf("service failed to create tenant: %w", err)
	}

	// Link creator to tenant
	if err := s.userService.AddTenantToUser(ctx, creatorID, tenant.ID); err != nil {
		// Optimization: Ideally we should rollback the created tenant here if this fails,
		// but for now we'll just return the error.
		return nil, fmt.Errorf("service failed to link creator to tenant: %w", err)
	}

	return tenant, nil
}
