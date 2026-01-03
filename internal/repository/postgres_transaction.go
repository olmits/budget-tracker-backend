package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/olmits/budget-tracker-backend/internal/models"
)

type PostgresTransactionRepo struct {
	DB *pgxpool.Pool
}

func (r *PostgresTransactionRepo) CreateTransaction(ctx context.Context, t *models.Transaction) error {
	sql := `INSERT INTO transactions (user_id, amount, description, date, category_id)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at`
	fmt.Println(t)

	return r.DB.QueryRow(ctx, sql,
		t.UserId, t.Amount, t.Description, t.Date, t.CategoryId,
	).Scan(&t.ID, &t.CreatedAt)
}

func (r *PostgresTransactionRepo) ListTransactions(ctx context.Context, userID uuid.UUID) ([]*models.Transaction, error) {
	// 1. Define the SQL
	// We use LEFT JOIN to fetch the Category Name if it exists
	sql := `SELECT
					t.id, t.amount, t.description, t.date, t.created_at,
					t.category_id, c.name as category_name
			FROM transactions t
			LEFT JOIN categories c ON t.category_id = c.id
			WHERE t.user_id = $1
			ORDER BY t.date DESC`
	// 2. Execute Query
	rows, err := r.DB.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 3. Iterate and Map to Structs
	var transactions []*models.Transaction

	for rows.Next() {
		t := &models.Transaction{}
		var catName *string // Handle NULL category name

		// Scan into the struct fields
		if err := rows.Scan(
			&t.ID,
			&t.Amount,
			&t.Description,
			&t.Date,
			&t.CreatedAt,
			&t.CategoryId,
			&catName,
		); err != nil {
			return nil, err
		}

		// Handle the pointer logic for category name
		if catName != nil {
			t.CategoryName = *catName
		}

		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
