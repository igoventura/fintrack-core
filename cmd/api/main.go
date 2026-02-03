package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/igoventura/fintrack-core/internal/api/handler"
	"github.com/igoventura/fintrack-core/internal/api/middleware"
	"github.com/igoventura/fintrack-core/internal/api/router"
	"github.com/igoventura/fintrack-core/internal/auth"
	"github.com/igoventura/fintrack-core/internal/db/postgres"
	"github.com/igoventura/fintrack-core/internal/service"
	"github.com/joho/godotenv"
)

// @title FinTrack API
// @version 1.0
// @description Financial tracking API documentation.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @SecurityDefinitions.oauth2.password AuthPassword
// @TokenUrl /auth/login
// @host localhost:8080
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	ctx := context.Background()

	// Database initialization
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db, err := postgres.NewDB(ctx, connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize Repositories
	accountRepo := postgres.NewAccountRepository(db)
	userRepo := postgres.NewUserRepository(db)
	tenantRepo := postgres.NewTenantRepository(db)

	// Initialize Services
	accountService := service.NewAccountService(accountRepo)
	userService := service.NewUserService(userRepo)

	// Initialize Handlers
	accountHandler := handler.NewAccountHandler(accountService)

	// Auth Middleware
	projectRef := os.Getenv("SUPABASE_PROJECT_REF")
	if projectRef == "" {
		log.Fatal("SUPABASE_PROJECT_REF environment variable is required")
	}

	// Construct JWKS URL: https://<project-ref>.supabase.co/auth/v1/.well-known/jwks.json
	jwksURL := "https://" + projectRef + ".supabase.co/auth/v1/.well-known/jwks.json"
	authValidator, err := auth.NewValidator(jwksURL)
	if err != nil {
		log.Fatalf("Failed to initialize auth validator: %v", err)
	}

	// Auth Service (Supabase)
	anonKey := os.Getenv("SUPABASE_ANON_KEY")
	if anonKey == "" {
		log.Fatal("SUPABASE_ANON_KEY environment variable is required")
	}
	authService := service.NewSupabaseAuthService(projectRef, anonKey, userService)
	authHandler := handler.NewAuthHandler(authService)

	// Create Middleware
	authMiddleware := middleware.NewAuthMiddleware(userRepo, authValidator)
	tenantMiddleware := middleware.NewTenantMiddleware(tenantRepo)

	// Router setup
	r := router.NewRouter(accountHandler, authHandler, authMiddleware, tenantMiddleware)

	// Server configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Server starting on port %s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %s", err)
	}
}
