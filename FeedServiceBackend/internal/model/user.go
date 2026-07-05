package models

import (
	"time"
)

// User представляет зарегистрированного пользователя.
//
// Содержит идентификатор, имя, email, дату создания и пароль.
// Пароль (Password) не сериализуется в JSON (json:"-") и предназначен
// только для внутренней обработки.
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"post_id"` // Примечание: тег json:"post_id" вероятно ошибочен, должно быть "created_at"
	Password  string    `json:"-"`
}
