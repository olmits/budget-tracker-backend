package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/olmits/budget-tracker-backend/internal/middleware"
	"github.com/olmits/budget-tracker-backend/internal/models"
	"github.com/olmits/budget-tracker-backend/internal/repository"
	"github.com/olmits/budget-tracker-backend/internal/service"
)

type CreateTransactionRequest struct {
	Amount      int64     `json:"amount" binding:"required"` // In cents!
	Description string    `json:"description"`
	Date        time.Time `json:"date" binding:"required"`
	CategoryId  *string   `json:"category_id"` // Optional
}

type TransactionHandler struct {
	Repo    repository.TransactionRepository
	Service *service.DashboardService
}

// POST /api/v1/transactions
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest

	// 1. Parse JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := middleware.GetUserID(c) // "9e0058fe-21e5-413b-bd89-bda904e9ba8d"
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
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
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         t.ID,
		"status":     "created",
		"created_at": t.CreatedAt,
	})
}

// GET /api/v1/transactions
func (h *TransactionHandler) ListTransactions(c *gin.Context) {
	userID, err := middleware.GetUserID(c) // "9e0058fe-21e5-413b-bd89-bda904e9ba8d"
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
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

// GET /api/v1/dashboard
func (h *TransactionHandler) GetDashboard(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	summary, err := h.Service.GetUserSummary(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate dashboard"})
		return
	}

	c.JSON(http.StatusOK, summary)
}
