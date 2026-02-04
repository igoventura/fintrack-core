package router

import (
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gin-gonic/gin"
	"github.com/igoventura/fintrack-api/internal/api/handler"
	"github.com/igoventura/fintrack-api/internal/api/middleware"
)

func NewRouter(accountHandler *handler.AccountHandler, authHandler *handler.AuthHandler, categoryHandler *handler.CategoryHandler, tagHandler *handler.TagHandler, tenantHandler *handler.TenantHandler, transactionHandler *handler.TransactionHandler, authMiddleware *middleware.AuthMiddleware, tenantMiddleware *middleware.TenantMiddleware, userHandler *handler.UserHandler) *gin.Engine {
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
	accounts.Use(authMiddleware.Handle(), tenantMiddleware.Handle(false))
	{
		accounts.GET("/", accountHandler.List)
		accounts.POST("/", accountHandler.Create)
		accounts.GET("/:id", accountHandler.Get)
	}

	// Category routes
	categories := r.Group("/categories")
	categories.Use(authMiddleware.Handle(), tenantMiddleware.Handle(false))
	{
		categories.GET("/", categoryHandler.ListCategories)
		categories.POST("/", categoryHandler.CreateCategory)
		categories.GET("/:id", categoryHandler.GetCategory)
		categories.PUT("/:id", categoryHandler.UpdateCategory)
		categories.DELETE("/:id", categoryHandler.DeleteCategory)
	}

	// Tag routes
	tags := r.Group("/tags")
	tags.Use(authMiddleware.Handle(), tenantMiddleware.Handle(false))
	{
		tags.GET("/", tagHandler.ListTags)
		tags.POST("/", tagHandler.CreateTag)
		tags.GET("/:id", tagHandler.GetTag)
		tags.PUT("/:id", tagHandler.UpdateTag)
		tags.DELETE("/:id", tagHandler.DeleteTag)
	}

	// Transaction routes
	transactions := r.Group("/transactions")
	transactions.Use(authMiddleware.Handle(), tenantMiddleware.Handle(false))
	{
		transactions.GET("/", transactionHandler.List)
		transactions.POST("/", transactionHandler.Create)
		transactions.GET("/:id", transactionHandler.GetByID)
		transactions.PUT("/:id", transactionHandler.Update)
		transactions.DELETE("/:id", transactionHandler.Delete)
	}

	// Auth routes
	auth := r.Group("/auth")
	auth.Use(tenantMiddleware.Handle(true))
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
