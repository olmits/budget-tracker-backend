package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/olmits/budget-tracker-backend/internal/models"
	"github.com/olmits/budget-tracker-backend/internal/repository"
)

type DashboardService struct {
	Repo repository.TransactionRepository
}

func (s *DashboardService) GetUserSummary(ctx context.Context, userID uuid.UUID) (*models.DashboardSummary, error) {
	sums, err := s.Repo.GetSummaryByType(ctx, userID)
	if err != nil {
		return nil, err
	}

	income := sums["income"]
	expense := sums["expense"]
	balance := income - expense

	return &models.DashboardSummary{
		TotalIncome:  income,
		TotalExpense: expense,
		NetBalance:   balance,
	}, nil
}
