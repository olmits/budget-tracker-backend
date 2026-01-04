package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/olmits/budget-tracker-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCategoryRepo struct {
	mock.Mock
}

func (m *MockCategoryRepo) CreateCategory(ctx context.Context, c *models.Category) error {
	args := m.Called(ctx, c)
	c.ID = uuid.New()
	return args.Error(0)
}

func (m *MockCategoryRepo) ListCategories(ctx context.Context, userID uuid.UUID) ([]*models.Category, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Category), args.Error(1)
}

func TestCreateCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Case", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		mockRepo.On("CreateCategory", mock.Anything, mock.Anything).Return(nil)

		dummyUserID := uuid.New()

		h := &CategoryHandler{Repo: mockRepo}
		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("userID", dummyUserID)
			ctx.Next()
		})
		r.POST("/api/v1/categories", h.CreateCategory)

		w := httptest.NewRecorder()
		jsonBody := []byte(`{"name": "Fun", "type": "expense"}`)
		req, _ := http.NewRequest("POST", "/api/v1/categories", bytes.NewBuffer(jsonBody))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "Fun")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Duplicate Error (Conflict)", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)

		pgErr := &pgconn.PgError{Code: "23505"}

		mockRepo.On("CreateCategory", mock.Anything, mock.Anything).Return(pgErr)

		dummyUserID := uuid.New()

		h := &CategoryHandler{Repo: mockRepo}
		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("userID", dummyUserID)
			ctx.Next()
		})
		r.POST("/api/v1/categories", h.CreateCategory)

		w := httptest.NewRecorder()
		jsonBody := []byte(`{"name": "Fun", "type": "expense"}`)
		req, _ := http.NewRequest("POST", "/api/v1/categories", bytes.NewBuffer(jsonBody))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "already exists")
		mockRepo.AssertExpectations(t)
	})
}

func TestListCategories(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Case", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		dummyList := []*models.Category{
			{Name: "Salary", Type: "income"},
			{Name: "Rent", Type: "expense"},
		}
		dummyUserID := uuid.New()

		mockRepo.On("ListCategories", mock.Anything, mock.Anything).Return(dummyList, nil)

		h := &CategoryHandler{Repo: mockRepo}
		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("userID", dummyUserID)
			ctx.Next()
		})
		r.GET("/api/v1/categories", h.ListCategories)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/categories", nil)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Salary")
		assert.Contains(t, w.Body.String(), "Rent")
	})
}
