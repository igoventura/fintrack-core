package service

import (
	"context"
	"fmt"

	"github.com/igoventura/fintrack-core/domain"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service failed to get user by id: %w", err)
	}
	return user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("service failed to get user by email: %w", err)
	}
	return user, nil
}

func (s *UserService) GetUserBySupabaseID(ctx context.Context, supabaseID string) (*domain.User, error) {
	user, err := s.repo.GetBySupabaseID(ctx, supabaseID)
	if err != nil {
		return nil, fmt.Errorf("service failed to get user by supabase id: %w", err)
	}
	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) error {
	if user.Email == "" {
		return fmt.Errorf("user email is required")
	}
	// Add other validation as needed

	if err := s.repo.Create(ctx, user); err != nil {
		return fmt.Errorf("service failed to create user: %w", err)
	}
	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *domain.User) error {
	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("service failed to update user: %w", err)
	}
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("service failed to delete user: %w", err)
	}
	return nil
}

func (s *UserService) AddTenantToUser(ctx context.Context, userID, tenantID string) error {
	if err := s.repo.AddUserToTenant(ctx, userID, tenantID); err != nil {
		return fmt.Errorf("service failed to add user to tenant: %w", err)
	}
	return nil
}

func (s *UserService) RemoveUserFromTenant(ctx context.Context, userID, tenantID string) error {
	if err := s.repo.RemoveUserFromTenant(ctx, userID, tenantID); err != nil {
		return fmt.Errorf("service failed to remove user from tenant: %w", err)
	}
	return nil
}

func (s *UserService) ListUserTenants(ctx context.Context, userID string) ([]domain.UserTenant, error) {
	tenants, err := s.repo.ListUserTenants(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service failed to list user tenants: %w", err)
	}
	return tenants, nil
}
