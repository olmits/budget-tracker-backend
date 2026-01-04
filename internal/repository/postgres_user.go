package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/olmits/budget-tracker-backend/internal/models"
)

type PostgresUserRepo struct {
	DB *pgxpool.Pool
}

func (r *PostgresUserRepo) CreateUser(ctx context.Context, user *models.User) error {
	sql := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at`
	return r.DB.QueryRow(ctx, sql, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)
}

func (r *PostgresUserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	sql := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`

	user := &models.User{}
	err := r.DB.QueryRow(ctx, sql, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}
