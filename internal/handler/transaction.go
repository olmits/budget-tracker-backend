package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateTransactionRequest struct {
	Amount      int64     `json:"amount" binding:"required"` // In cents!
	Description string    `json:"description"`
	Date        time.Time `json:"date" binding:"required"`
	CategoryId  *string   `json:"category_id"` // Optional
}

type TransactionHandler struct {
	DB *pgxpool.Pool
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest

	// 1. Parse JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. TODO: In the future, get this ID from the Auth Token (JWT)
	// For now, HARDCODE the ID you got from the database step above!
	userID := "9e0058fe-21e5-413b-bd89-bda904e9ba8d"

	// 3. Insert into Database
	sql := `INSERT INTO transactions (user_id, amount, description, date, category_id)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at`

	var id string
	var createdAt time.Time

	// Note: We handle the optional CategoryID logic here
	err := h.DB.QueryRow(c.Request.Context(), sql,
		userID,
		req.Amount,
		req.Description,
		req.Date,
		req.CategoryId, // pgx handles nil pointers automatically for SQL NULL
	).Scan(&id, &createdAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction: " + err.Error()})
	}

	// 4. Return Success Response
	c.JSON(http.StatusCreated, gin.H{
		"id":         id,
		"status":     "created",
		"created_at": createdAt,
	})
}
