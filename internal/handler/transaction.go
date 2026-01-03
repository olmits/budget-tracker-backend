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

func (h *TransactionHandler) ListTransactions(c *gin.Context) {
	// 1. Hardcode UserID (Until we add Auth)
	userID := "9e0058fe-21e5-413b-bd89-bda904e9ba8d"

	// 2. Define the SQL
	// We use LEFT JOIN to fetch the Category Name if it exists
	sql := `SELECT
					t.id, t.amount, t.description, t.date, t.category_id,
					c.name as category_name
			FROM transactions t
			LEFT JOIN categories c ON t.category_id = c.id
			WHERE t.user_id = $1
			ORDER BY t.date DESC`

	// 3. Run the Query
	rows, err := h.DB.Query(c.Request.Context(), sql, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}
	defer rows.Close() // Important: Close the connection when done!

	// 4. Iterate over the results
	// We create a slice (array) to hold the data

	results := []map[string]interface{}{}

	for rows.Next() {
		var id string
		var amount int64
		var description string
		var date time.Time
		var categoryId *string
		var categoryName *string // Use pointer because category might be null

		// Scan the columns into variables
		if err := rows.Scan(&id, &amount, &description, &date, &categoryId, &categoryName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row"})
			return
		}

		// Append to our list
		results = append(results, map[string]interface{}{
			"id":            id,
			"amount":        amount,
			"description":   description,
			"date":          date,
			"category_id":   categoryId,
			"category_name": categoryName,
		})
	}

	// 5. Check for errors that happened during iteration
	if rows.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Row iteration error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}
