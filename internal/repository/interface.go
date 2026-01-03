package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/olmits/budget-tracker-backend/internal/models"
)

// TransactionRepository defines "what" we need from the DB, not "how"
type TransactionRepository interface {
	CreateTransaction(ctx context.Context, t *models.Transaction) error
	ListTransactions(ctx context.Context, userID uuid.UUID) ([]*models.Transaction, error)
}

type CategoryRepository interface {
	CreateCategory(ctx context.Context, c *models.Category) error
	ListCategories(ctx context.Context, userID uuid.UUID) ([]*models.Category, error)
}
