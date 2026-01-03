package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryHandler struct {
	DB *pgxpool.Pool
}

// Define the input JSON structure
type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required,oneof=income expense"` // Validation!
}

// POST /api/v1/categories
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest

	// 1. Validate Input (Gin handles the "oneof" check automatically)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Hardcode User ID (Temporary)
	userID := "9e0058fe-21e5-413b-bd89-bda904e9ba8d"

	// 3. Insert into DB
	sql := `INSERT INTO categories (user_id, name, type)
			VALUES ($1, $2, $3)
			RETURNING id, created_at`

	var id string
	var createdAt time.Time

	// Execute Query
	err := h.DB.QueryRow(c.Request.Context(), sql, userID, req.Name, req.Type).Scan(&id, &createdAt)

	if err != nil {
		// 1. Check if the error is a Postgres Error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// 2. Check for Error Code "23505" (Unique Violation)
			if pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, gin.H{
					"error": "Category with this name already exists",
				})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         id,
		"name":       req.Name,
		"type":       req.Type,
		"created_at": createdAt,
	})
}

// GET /api/v1/categories
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	userID := "9e0058fe-21e5-413b-bd89-bda904e9ba8d"

	sql := `SELECT id, name, type FROM categories WHERE user_id = $1 ORDER BY name ASC`

	rows, err := h.DB.Query(c.Request.Context(), sql, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	defer rows.Close()

	results := []map[string]string{}

	// 1. Iterate
	for rows.Next() {
		var id, name, catType string
		if err := rows.Scan(&id, &name, &catType); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scanning error"})
			return
		}
		results = append(results, map[string]string{
			"id":   id,
			"name": name,
			"type": catType,
		})
	}

	// 2. IMPORTANT: Check if the loop stopped because of an error
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database iteration error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}
