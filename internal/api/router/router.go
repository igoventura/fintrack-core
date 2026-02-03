package router

import (
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gin-gonic/gin"
	"github.com/igoventura/fintrack-core/internal/api/handler"
	"github.com/igoventura/fintrack-core/internal/api/middleware"
)

func NewRouter(accountHandler *handler.AccountHandler, authHandler *handler.AuthHandler, tenantHandler *handler.TenantHandler, authMiddleware *middleware.AuthMiddleware, tenantMiddleware *middleware.TenantMiddleware, userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// Tenant routes
	tenants := r.Group("/tenants")
	tenants.Use(authMiddleware.Handle())
	{
		tenants.POST("/", tenantHandler.Create)
	}

	// Account routes
	accounts := r.Group("/accounts")
	accounts.Use(authMiddleware.Handle(), tenantMiddleware.Handle())
	{
		accounts.GET("/", accountHandler.List)
		accounts.POST("/", accountHandler.Create)
		accounts.GET("/:id", accountHandler.Get)
	}

	// Auth routes
	auth := r.Group("/auth")
	auth.Use(tenantMiddleware.Handle())
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// User routes
	users := r.Group("/users")
	users.Use(authMiddleware.Handle())
	{
		users.GET("/profile", userHandler.GetProfile)
		users.PUT("/profile", userHandler.UpdateProfile)
	}

	// Documentation
	r.GET("/docs", func(c *gin.Context) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/swagger.yaml",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "FinTrack API",
			},
			DarkMode: true,
		})
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, htmlContent)
	})

	r.StaticFile("swagger.yaml", "./docs/swagger.yaml")

	return r
}
