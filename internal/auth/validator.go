package auth

import (
	"fmt"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

// SupabaseClaims represents the claims in a Supabase JWT
type SupabaseClaims struct {
	jwt.RegisteredClaims
	Email        string                 `json:"email"`
	Role         string                 `json:"role"`
	AppMetadata  map[string]interface{} `json:"app_metadata"`
	UserMetadata map[string]interface{} `json:"user_metadata"`
}

// Validator handles JWT validation using JWKS
type Validator struct {
	jwks keyfunc.Keyfunc
}

// NewValidator creates a new Validator with the given JWKS URL
func NewValidator(jwksURL string) (*Validator, error) {
	// Create the JWKS from the given URL.
	jwks, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("failed to create JWKS from resource at %s: %w", jwksURL, err)
	}

	return &Validator{jwks: jwks}, nil
}

// ValidateToken verifies the Supabase JWT signature and expiry using JWKS.
func (v *Validator) ValidateToken(tokenString string) (*SupabaseClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &SupabaseClaims{}, v.jwks.Keyfunc)

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(*SupabaseClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}
