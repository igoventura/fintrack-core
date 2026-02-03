package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/igoventura/fintrack-core/internal/api/handler"
	"github.com/igoventura/fintrack-core/internal/api/middleware"
	"github.com/mvrilo/go-redoc"
)

func NewRouter(accountHandler *handler.AccountHandler, authMiddleware *middleware.AuthMiddleware, tenantMiddleware *middleware.TenantMiddleware) *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// Account routes
	accounts := r.Group("/accounts")
	accounts.Use(authMiddleware.Handle(), tenantMiddleware.Handle())
	{
		accounts.GET("/", accountHandler.List)
		accounts.POST("/", accountHandler.Create)
		accounts.GET("/:id", accountHandler.Get)
	}

	// Documentation
	doc := redoc.Redoc{
		Title:       "FinTrack API",
		Description: "Financial tracking API documentation",
		SpecFile:    "./docs/swagger.yaml",
		SpecPath:    "/swagger.yaml",
		DocsPath:    "/docs",
	}

	// Convert redoc handler to gin compatible if needed,
	// but Redoc.Handler() returns a standard http.Handler.
	// We can use gin.WrapH to wrap it.
	r.GET("/docs", gin.WrapH(doc.Handler()))
	r.StaticFile("/swagger.yaml", "./docs/swagger.yaml")

	return r
}
