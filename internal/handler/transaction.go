package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/olmits/budget-tracker-backend/internal/models"
	"github.com/olmits/budget-tracker-backend/internal/repository"
)

type CreateTransactionRequest struct {
	Amount      int64     `json:"amount" binding:"required"` // In cents!
	Description string    `json:"description"`
	Date        time.Time `json:"date" binding:"required"`
	CategoryId  *string   `json:"category_id"` // Optional
}

type TransactionHandler struct {
	Repo repository.TransactionRepository
}

// POST /api/v1/transactions
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest

	// 1. Parse JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. TODO: In the future, get this ID from the Auth Token (JWT)
	// For now, HARDCODE the ID you got from the database step above!
	rawUserID := "9e0058fe-21e5-413b-bd89-bda904e9ba8d"
	userID, err := uuid.Parse(rawUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID format"})
		return
	}

	// Map Request to Model
	t := &models.Transaction{
		UserId:      userID,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        req.Date,
	}

	// CALL THE INTERFACE
	if err := h.Repo.CreateTransaction(c.Request.Context(), t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
	}

	// Return Success Response
	c.JSON(http.StatusCreated, gin.H{
		"id":         t.ID,
		"status":     "created",
		"created_at": t.CreatedAt,
	})
}

// GET /api/v1/transactions
func (h *TransactionHandler) ListTransactions(c *gin.Context) {
	// 1. Hardcode UserID (Until we add Auth)
	// We need to parse the string UUID into a real UUID object
	rawUserID := "9e0058fe-21e5-413b-bd89-bda904e9ba8d"
	userID, err := uuid.Parse(rawUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID format"})
		return
	}

	// 2. Call the Repository
	transactions, err := h.Repo.ListTransactions(c.Request.Context(), userID)
	if err != nil {
		// Log the error internally here if you have a logger
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	// 3. Return JSON
	c.JSON(http.StatusOK, gin.H{"data": transactions})
}
