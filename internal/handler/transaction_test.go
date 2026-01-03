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

// 1. Define the Mock Object
type MockTransactionRepo struct {
	mock.Mock
}

// Implement the Interface method for the Mock
func (m *MockTransactionRepo) CreateTransaction(ctx context.Context, t *models.Transaction) error {
	// This records that the method was called
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *MockTransactionRepo) ListTransactions(ctx context.Context, userID uuid.UUID) ([]*models.Transaction, error) {
	args := m.Called(ctx, userID)
	// We must cast the first return value to the specific type
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Transaction), args.Error(1)
}

func TestCreateTransaction(t *testing.T) {
	// Setup Gin to Test Mode (quieter logs)
	gin.SetMode(gin.TestMode)

	t.Run("Success Case", func(t *testing.T) {
		// A. Arrange (Setup)
		mockRepo := new(MockTransactionRepo)

		// Expect CreateTransaction to be called once, return NO error (nil)
		mockRepo.On("CreateTransaction", mock.Anything, mock.Anything).Return(nil)

		// Create Handler with the Mock
		h := &TransactionHandler{Repo: mockRepo}

		// Create a Request
		w := httptest.NewRecorder()
		jsonBody := []byte(`{"amount": 1000, "date": "2023-10-27T10:00:00Z", "description": "Test"}`)
		req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(jsonBody))

		// Setup Router
		r := gin.Default()
		r.POST("/api/v1/transactions", h.CreateTransaction)

		// B. Act (Run)
		r.ServeHTTP(w, req)

		// C. Assert (Check results)
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "created")

		mockRepo.AssertExpectations(t)
	})

	t.Run("Database Error Case", func(t *testing.T) {
		// A. Arrange (Setup)
		mockRepo := new(MockTransactionRepo)

		// Expect it to be called, but return an error this time
		mockRepo.On("CreateTransaction", mock.Anything, mock.Anything).Return(errors.New("db connection lost"))

		// Create Handler with the Mock
		h := &TransactionHandler{Repo: mockRepo}

		// Create a Request
		w := httptest.NewRecorder()
		jsonBody := []byte(`{"amount": 1000, "date": "2023-10-27T10:00:00Z", "description": "Test"}`)
		req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(jsonBody))

		// Setup Router
		r := gin.Default()
		r.POST("/api/v1/transactions", h.CreateTransaction)

		// B. Act (Run)
		r.ServeHTTP(w, req)

		// C. Assert (Check results)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to create")
	})
}

func TestListTransactions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Case", func(t *testing.T) {
		// A. Arrange
		mockRepo := new(MockTransactionRepo)

		// Create some dummy data to return
		dummyID := uuid.New()
		dummyTransactions := []*models.Transaction{
			{
				ID:           dummyID,
				Amount:       5000,
				Description:  "Test desc",
				CategoryName: "Test cat",
			},
		}

		// Expect ListTransactions to be called with ANY UserID
		mockRepo.On("ListTransactions", mock.Anything, mock.Anything).Return(dummyTransactions, nil)

		// Create Handler with the Mock
		h := &TransactionHandler{Repo: mockRepo}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/transactions", nil)

		r := gin.Default()
		r.GET("/api/v1/transactions", h.ListTransactions)

		// B. Act
		r.ServeHTTP(w, req)

		// C. Assert
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify the JSON body contains our data
		assert.Contains(t, w.Body.String(), "Test desc")
		assert.Contains(t, w.Body.String(), "5000")
		assert.Contains(t, w.Body.String(), "Test cat")

		mockRepo.AssertExpectations(t)
	})

	t.Run("Database Error Case", func(t *testing.T) {
		// A. Arrange
		mockRepo := new(MockTransactionRepo)

		// Expect it to fail
		mockRepo.On("ListTransactions", mock.Anything, mock.Anything).Return(nil, errors.New("db disconnected"))

		h := &TransactionHandler{Repo: mockRepo}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/transactions", nil)

		r := gin.Default()
		r.GET("/api/v1/transactions", h.ListTransactions)

		// B. Act
		r.ServeHTTP(w, req)

		// C. Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to fetch transactions")

		mockRepo.AssertExpectations(t)
	})
}
