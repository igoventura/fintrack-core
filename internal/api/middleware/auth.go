package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/igoventura/fintrack-api/domain"
	"github.com/igoventura/fintrack-api/internal/auth"
)

const (
	UserIDKey = "userID"
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

		user, err := m.userRepo.GetBySupabaseID(c.Request.Context(), claims.Subject)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		c.Set(UserIDKey, user.ID)
		ctx := domain.WithUserID(c.Request.Context(), user.ID)
		ctx = domain.WithToken(ctx, tokenString)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
