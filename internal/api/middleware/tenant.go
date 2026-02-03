package middleware

import (
	"github.com/gin-gonic/gin"
)

const (
	TenantIDKey    = "tenantID"
	TenantIDHeader = "X-Tenant-ID"
)

type TenantMiddleware struct{}

func NewTenantMiddleware() *TenantMiddleware {
	return &TenantMiddleware{}
}

func (m *TenantMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader(TenantIDHeader)
		if tenantID != "" {
			c.Set(TenantIDKey, tenantID)
		}
		c.Next()
	}
}
