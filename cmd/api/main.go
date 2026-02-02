package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/igoventura/fintrack-core/internal/api/handler"
	"github.com/igoventura/fintrack-core/internal/api/router"
	"github.com/igoventura/fintrack-core/internal/db/postgres"
	"github.com/igoventura/fintrack-core/internal/service"
	"github.com/joho/godotenv"
)

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

	// Initialize Services
	accountService := service.NewAccountService(accountRepo)

	// Initialize Handlers
	accountHandler := handler.NewAccountHandler(accountService)

	// Router setup
	r := router.NewRouter(accountHandler)

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
