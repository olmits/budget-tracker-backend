package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never send password in JSON response
	CreatedAt    time.Time `json:"created_at"`
}

type Category struct {
	ID        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // "income" or "expense"
	CreatedAt time.Time `json:"created_at"`
}

type Transaction struct {
	ID          uuid.UUID  `json:"id"`
	UserId      uuid.UUID  `json:"user_id"`
	CategoryId  *uuid.UUID `json:"category_id"` // Pointer because it can be null
	Amount      int64      `json:"amount"`      // Cents
	Description string     `json:"description"`
	Date        time.Time  `json:"date"`
	CreatedAt   time.Time  `json:"created_at"`
}
