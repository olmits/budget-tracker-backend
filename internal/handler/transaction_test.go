package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/olmits/budget-tracker-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 1. Mock the Repository (Because your handler calls h.Repo)
type MockTransactionRepo struct {
	mock.Mock
}

func (m *MockTransactionRepo) CreateTransaction(ctx context.Context, t *models.Transaction) error {
	// mock.Anything is used for context, mock.AnythingOfType for the transaction struct
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *MockTransactionRepo) ListTransactions(ctx context.Context, userID uuid.UUID) ([]*models.Transaction, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepo) GetSummaryByType(ctx context.Context, userID uuid.UUID) (map[string]int64, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int64), args.Error(1)
}

func TestCreateTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// A valid UUID required by your new validation rules
	validCategoryID := "550e8400-e29b-41d4-a716-446655440000"
	dummyUserID := uuid.New()

	t.Run("Success Case", func(t *testing.T) {
		// A. Arrange
		mockRepo := new(MockTransactionRepo)

		// Expect Repo to be called with a context and a *models.Transaction
		mockRepo.On("CreateTransaction", mock.Anything, mock.AnythingOfType("*models.Transaction")).Return(nil)

		// Create Handler using the Mock Repo
		h := &TransactionHandler{
			Repo: mockRepo,
			// Service is nil because CreateTransaction doesn't use it
		}

		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			// Assuming your middleware.GetUserID expects a uuid.UUID in context
			ctx.Set("userID", dummyUserID)
			ctx.Next()
		})
		r.POST("/api/v1/transactions", h.CreateTransaction)

		w := httptest.NewRecorder()

		// B. Act
		// FIX: We added "category_id" to match your binding:"required"
		jsonBody := []byte(`{
			"amount": 1000, 
			"date": "2023-10-27T10:00:00Z", 
			"description": "Test",
			"category_id": "` + validCategoryID + `"
		}`)
		req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(jsonBody))

		r.ServeHTTP(w, req)

		// C. Assert
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "created")

		mockRepo.AssertExpectations(t)
	})

	t.Run("Database Error Case", func(t *testing.T) {
		// A. Arrange
		mockRepo := new(MockTransactionRepo)

		// Expect it to fail this time
		mockRepo.On("CreateTransaction", mock.Anything, mock.AnythingOfType("*models.Transaction")).Return(errors.New("db error"))

		h := &TransactionHandler{Repo: mockRepo}

		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("userID", dummyUserID)
			ctx.Next()
		})
		r.POST("/api/v1/transactions", h.CreateTransaction)

		w := httptest.NewRecorder()
		jsonBody := []byte(`{
			"amount": 1000, 
			"date": "2023-10-27T10:00:00Z", 
			"description": "Test",
			"category_id": "` + validCategoryID + `"
		}`)
		req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(jsonBody))

		// B. Act
		r.ServeHTTP(w, req)

		// C. Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to create transaction")

		mockRepo.AssertExpectations(t)
	})
}
