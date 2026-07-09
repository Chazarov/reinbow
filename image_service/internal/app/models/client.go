package models

import (
	"time"
)

// Client представляет учётную запись клиента сервиса.
type Client struct {
	Id        int       `json:"id" db:"id"`
	PassHash  string    `json:"pass_hash" db:"pass_hash"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CreateClientMto передаёт данные для создания клиента в репозитории.
type CreateClientMto struct {
	PassHash  string    `json:"pass_hash" db:"pass_hash"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CreateClientWithPasswordMto передаёт данные регистрации клиента с открытым паролем.
type CreateClientWithPasswordMto struct {
	Password  string    `json:"password"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
