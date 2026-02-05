package service

import (
	"context"

	"github.com/igoventura/fintrack-api/domain"
	"github.com/igoventura/fintrack-api/internal/api/dto"
	"github.com/supabase-community/gotrue-go"
	"github.com/supabase-community/gotrue-go/types"
)

type AuthService interface {
	Register(ctx context.Context, email, password, fullName string) (*dto.AuthResponse, error)
	Login(ctx context.Context, email, password string) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.AuthResponse, error)
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

func (s *SupabaseAuthService) Register(ctx context.Context, email, password, fullName string) (*dto.AuthResponse, error) {
	req := types.SignupRequest{
		Email:    email,
		Password: password,
		Data: map[string]any{
			"full_name": fullName,
		},
	}
	resp, err := s.client.Signup(req)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:      email,
		SupabaseID: resp.User.ID.String(),
		Name:       fullName,
	}
	if err := s.userService.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	if tenantID := domain.GetTenantID(ctx); tenantID != "" {
		if err := s.userService.AddTenantToUser(ctx, user.ID, tenantID); err != nil {
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

func (s *SupabaseAuthService) UpdateUser(ctx context.Context, user *domain.User) error {
	token := domain.GetToken(ctx)
	updateUserRequest := types.UpdateUserRequest{
		Email: user.Email,
		Data: map[string]any{
			"full_name": user.Name,
		},
	}
	if _, err := s.client.WithToken(token).UpdateUser(updateUserRequest); err != nil {
		return err
	}

	return nil
}

func (s *SupabaseAuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.AuthResponse, error) {
	resp, err := s.client.RefreshToken(refreshToken)
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
