package service

import (
	"context"

	"github.com/igoventura/fintrack-core/domain"
	"github.com/igoventura/fintrack-core/internal/api/dto"
	"github.com/supabase-community/gotrue-go"
	"github.com/supabase-community/gotrue-go/types"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (*dto.AuthResponse, error)
	Login(ctx context.Context, email, password string) (*dto.AuthResponse, error)
}

type SupabaseAuthService struct {
	client      gotrue.Client
	userService *UserService
}

func NewSupabaseAuthService(projectID, apiKey string, userService *UserService) *SupabaseAuthService {
	return &SupabaseAuthService{
		client:      gotrue.New(projectID, apiKey),
		userService: userService,
	}
}

func (s *SupabaseAuthService) Register(ctx context.Context, email, password string) (*dto.AuthResponse, error) {
	req := types.SignupRequest{
		Email:    email,
		Password: password,
	}
	resp, err := s.client.Signup(req)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:      email,
		SupabaseID: resp.User.ID.String(),
		Name:       resp.User.UserMetadata["full_name"].(string),
	}
	if err := s.userService.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	if tenantID := domain.GetTenantID(ctx); tenantID != "" {
		if err := s.userService.AddUserToTenant(ctx, user.ID, tenantID); err != nil {
			return nil, err
		}
	}

	return &dto.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		User: dto.User{
			ID:    resp.User.ID.String(),
			Email: resp.User.Email,
		},
	}, nil
}

func (s *SupabaseAuthService) Login(ctx context.Context, email, password string) (*dto.AuthResponse, error) {
	resp, err := s.client.SignInWithEmailPassword(email, password)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		User: dto.User{
			ID:    resp.User.ID.String(),
			Email: resp.User.Email,
		},
	}, nil
}
