package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/igoventura/fintrack-core/domain"
)

const (
	TenantIDHeader = "X-Tenant-ID"
)

type TenantMiddleware struct {
	tenantRepo domain.TenantRepository
}

func NewTenantMiddleware(tenantRepo domain.TenantRepository) *TenantMiddleware {
	return &TenantMiddleware{tenantRepo: tenantRepo}
}

func (m *TenantMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader(TenantIDHeader)
		if tenantID != "" {
			if _, err := m.tenantRepo.GetByID(c.Request.Context(), tenantID); err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "tenant ID is not valid"})
				return
			}
			ctx := domain.WithTenantID(c.Request.Context(), tenantID)
			c.Request = c.Request.WithContext(ctx)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "tenant ID is required"})
			return
		}
		c.Next()
	}
}
