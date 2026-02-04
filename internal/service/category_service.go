package service

import (
	"context"
	"fmt"

	"github.com/igoventura/fintrack-core/domain"
)

type CategoryService struct {
	repo domain.CategoryRepository
}

func NewCategoryService(repo domain.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetCategory(ctx context.Context, id string) (*domain.Category, error) {
	tenantID := domain.GetTenantID(ctx)
	category, err := s.repo.GetByID(ctx, id, tenantID)
	if err != nil {
		return nil, fmt.Errorf("service failed to get category: %w", err)
	}
	return category, nil
}

func (s *CategoryService) ListCategories(ctx context.Context) ([]domain.Category, error) {
	tenantID := domain.GetTenantID(ctx)
	categories, err := s.repo.List(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("service failed to list categories: %w", err)
	}
	return categories, nil
}

func (s *CategoryService) CreateCategory(ctx context.Context, category *domain.Category) error {
	tenantID := domain.GetTenantID(ctx)
	category.TenantID = tenantID

	isValid, validationErrors := category.IsValid()
	if !isValid {
		return fmt.Errorf("invalid category: %v", validationErrors)
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return fmt.Errorf("service failed to create category: %w", err)
	}
	return nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, category *domain.Category) error {
	tenantID := domain.GetTenantID(ctx)
	category.TenantID = tenantID

	// 1. Validate input
	isValid, validationErrors := category.IsValid()
	if !isValid {
		return fmt.Errorf("invalid category: %v", validationErrors)
	}

	// 2. Check existence and ownership (Tenant Isolation)
	// We just check if it exists for this tenant
	if _, err := s.repo.GetByID(ctx, category.ID, tenantID); err != nil {
		return fmt.Errorf("service failed to get category for update (or unauthorized): %w", err)
	}

	// 3. Update
	if err := s.repo.Update(ctx, category); err != nil {
		return fmt.Errorf("service failed to update category: %w", err)
	}
	return nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id, userID string) error {
	tenantID := domain.GetTenantID(ctx)

	// 1. Check existence and ownership
	if _, err := s.repo.GetByID(ctx, id, tenantID); err != nil {
		return fmt.Errorf("service failed to get category for delete (or unauthorized): %w", err)
	}

	// 2. Delete
	if err := s.repo.Delete(ctx, id, tenantID, userID); err != nil {
		return fmt.Errorf("service failed to delete category: %w", err)
	}
	return nil
}
