package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/olmits/budget-tracker-backend/internal/models"
)

type PostgresCategoryRepo struct {
	DB *pgxpool.Pool
}

func (r *PostgresCategoryRepo) CreateCategory(ctx context.Context, c *models.Category) error {
	sql := `INSERT INTO categories (user_id, name, type)
			VALUES ($1, $2, $3)
			RETURNING id, created_at`
	return r.DB.QueryRow(ctx, sql,
		c.UserId, c.Name, c.Type,
	).Scan(&c.ID, &c.CreatedAt)
}

func (r *PostgresCategoryRepo) ListCategories(ctx context.Context, userID uuid.UUID) ([]*models.Category, error) {
	// 1. Define the SQL
	sql := `SELECT id, name, type FROM categories WHERE user_id = $1 ORDER BY name ASC`

	// 2. Execute Query
	rows, err := r.DB.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 3. Iterate and Map to Structs
	var categories []*models.Category

	for rows.Next() {
		c := &models.Category{}
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Type,
		); err != nil {
			return nil, err
		}
		// Since we just scanned name/type, we should also manually set the UserID
		c.UserId = userID
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
