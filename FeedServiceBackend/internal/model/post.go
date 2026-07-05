package models

// Post представляет публикацию пользователя в ленте.
//
// Содержит основную информацию: идентификатор, автора, заголовок, описание,
// дату создания и список прикреплённых изображений (ID). Поля Username и Images
// опциональны (omitempty) и заполняются при необходимости, например, при
// выдаче ленты.
type Post struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Username    string `json:"username,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	Images      []int  `json:"images,omitempty"`
}
