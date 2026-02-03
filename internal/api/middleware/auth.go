package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/igoventura/fintrack-core/domain"
	"github.com/igoventura/fintrack-core/internal/auth"
)

const (
	UserContextKey = "currentUser"
)

type AuthMiddleware struct {
	userRepo  domain.UserRepository
	validator *auth.Validator
}

func NewAuthMiddleware(userRepo domain.UserRepository, validator *auth.Validator) *AuthMiddleware {
	return &AuthMiddleware{
		userRepo:  userRepo,
		validator: validator,
	}
}

func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		tokenString := parts[1]
		claims, err := m.validator.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Check if user exists in DB by Supabase ID (sub claim)
		user, err := m.userRepo.GetBySupabaseID(c.Request.Context(), claims.Subject)
		if err != nil {
			// If user doesn't exist, we might want to create it (Auto-provisioning)
			// For now, let's assume strict checking, or we can implement auto-creation here if requested.
			// Since we want to link users, if not found, it implies they haven't been synced or created yet.
			// To be safe, we'll return Unauthorized if not found, or maybe Forbidden?
			// But clean arch suggests maybe `Handle` shouldn't do complex creating logic?
			// A simple approach: return 401 if user not found in our system.
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		c.Set(UserContextKey, user)
		c.Next()
	}
}
