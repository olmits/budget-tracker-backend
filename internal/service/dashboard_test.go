package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/olmits/budget-tracker-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) GetSummaryByType(ctx context.Context, userID uuid.UUID) (map[string]int64, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *MockRepo) CreateTransaction(ctx context.Context, t *models.Transaction) error {
	return nil // Not used in this test
}
func (m *MockRepo) ListTransactions(ctx context.Context, userID uuid.UUID) ([]*models.Transaction, error) {
	return nil, nil // Not used in this test
}

func TestGetUserDashboard(t *testing.T) {
	t.Run("Calculates Balance Correctly", func(t *testing.T) {
		mockRepo := new(MockRepo)
		userID := uuid.New()

		mockData := map[string]int64{
			"income":  1000,
			"expense": 300,
		}

		mockRepo.On("GetSummaryByType", mock.Anything, userID).Return(mockData, nil)

		s := &DashboardService{Repo: mockRepo}

		summary, err := s.GetUserSummary(context.Background(), userID)

		assert.NoError(t, err)
		assert.NotNil(t, summary)

		assert.Equal(t, int64(1000), summary.TotalIncome)
		assert.Equal(t, int64(300), summary.TotalExpense)
		assert.Equal(t, int64(700), summary.NetBalance)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Handles DB Erorr", func(t *testing.T) {
		mockRepo := new(MockRepo)
		userID := uuid.New()

		mockRepo.On("GetSummaryByType", mock.Anything, userID).Return(nil, errors.New("db disconnected"))

		s := &DashboardService{Repo: mockRepo}

		summary, err := s.GetUserSummary(context.Background(), userID)

		assert.Error(t, err)
		assert.Nil(t, summary)

		assert.Equal(t, "db disconnected", err.Error())

		mockRepo.AssertExpectations(t)
	})
}
