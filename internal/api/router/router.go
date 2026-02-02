package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/igoventura/fintrack-core/internal/api/handler"
	"github.com/mvrilo/go-redoc"
)

func NewRouter(accountHandler *handler.AccountHandler) *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// Account routes
	accounts := r.Group("/accounts")
	{
		accounts.GET("/", accountHandler.List)
		accounts.POST("/", accountHandler.Create)
		accounts.GET("/:id", accountHandler.Get)
	}

	// Documentation
	doc := redoc.Redoc{
		Title:       "FinTrack API",
		Description: "Financial tracking API documentation",
		SpecFile:    "./docs/openapi.yaml",
		SpecPath:    "/openapi.yaml",
		DocsPath:    "/docs",
	}

	// Convert redoc handler to gin compatible if needed,
	// but Redoc.Handler() returns a standard http.Handler.
	// We can use gin.WrapH to wrap it.
	r.GET("/docs", gin.WrapH(doc.Handler()))
	r.StaticFile("/openapi.yaml", "./docs/openapi.yaml")

	return r
}
