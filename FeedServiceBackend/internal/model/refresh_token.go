package models

import (
	"time"
)

// RefreshToken представляет запись о refresh-токене в базе данных.
// Содержит информацию о пользователе, самом токене и сроках его действия.
type RefreshToken struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
