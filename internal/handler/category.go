package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/olmits/budget-tracker-backend/internal/models"
	"github.com/olmits/budget-tracker-backend/internal/repository"
)

type CategoryHandler struct {
	Repo repository.CategoryRepository
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

	// 2. TODO: In the future, get this ID from the Auth Token (JWT)
	// For now, HARDCODE the ID you got from the database step above!
	rawUserID := "9e0058fe-21e5-413b-bd89-bda904e9ba8d"
	userID, err := uuid.Parse(rawUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID format"})
		return
	}

	cat := &models.Category{
		UserId: userID,
		Name:   req.Name,
		Type:   req.Type,
	}

	// 3. Call repository
	if err := h.Repo.CreateCategory(c.Request.Context(), cat); err != nil {
		// Check if the error is a Postgres Error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Category with this name already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, cat)
}

// GET /api/v1/categories
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	// 1. Hardcode UserID (Until we add Auth)
	// We need to parse the string UUID into a real UUID object
	rawUserID := "9e0058fe-21e5-413b-bd89-bda904e9ba8d"
	userID, err := uuid.Parse(rawUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID format"})
		return
	}

	// 2. Call the Repository
	categories, err := h.Repo.ListCategories(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": categories})
}
