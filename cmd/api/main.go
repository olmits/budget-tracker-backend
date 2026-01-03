package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/olmits/budget-tracker-backend/internal/handler"
	"github.com/olmits/budget-tracker-backend/internal/repository"
	"github.com/olmits/budget-tracker-backend/pkg/database"
)

func main() {
	// 1. Define Database Credentials (typically loaded from .env)
	// Format: postgres://user:password@host:port/dbname
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// 2. Connect to the Database
	fmt.Println("Connecting to database...")
	dbPool, err := database.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close() // Close connection when main() exits

	fmt.Println("Successfully connected to PostgreSQL!")

	// 3. Initialize the Repository layer
	// We pass the DB pool into the concrete Postgres implementation
	transactionRepo := &repository.PostgresTransactionRepo{DB: dbPool}
	categoryRepo := &repository.PostgresCategoryRepo{DB: dbPool}

	// 4. Initialize the Handler layer
	txHandler := &handler.TransactionHandler{Repo: transactionRepo}
	catHandler := &handler.CategoryHandler{Repo: categoryRepo}

	// 5. Initialize the Router (Gin)
	r := gin.Default()

	// 6. Simple Health Check Endpoint (using the DB)
	r.GET("/health", func(c *gin.Context) {
		// Ping the DB again to verify it's still alive
		if err := dbPool.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Database disconnected"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "active", "database": "connected"})
	})

	// Transaction Routes
	r.POST("/api/v1/transactions", txHandler.CreateTransaction)
	r.GET("/api/v1/transactions", txHandler.ListTransactions)

	// Category Routes
	r.POST("/api/v1/categories", catHandler.CreateCategory)
	r.GET("/api/v1/categories", catHandler.ListCategories)

	// 7. Start Server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = ":8080"
	}

	fmt.Printf("Starting server on port %s...\n", port)
	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}
