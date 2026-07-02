package models

import (
	"time"
)

type Post struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Username    string    `json:"username,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   string	  `json:"created_at"`
	Images      []int  	  `json:"images,omitempty"`
}

type ImagePost struct {
	ID     	int `json:"id"`
	PostID 	int `json:"post_id"`
	ImageID int `json:"image_id"`
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email  	  string    `json:"email"`
	TgChatId string     `json:"tg_chat_id"`
	CreatedAt time.Time `json:"post_id"`
	Password  string    `json:"-"`
}

type RefreshToken struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    Token     string    `json:"token"`
    ExpiresAt time.Time `json:"expires_at"`
    CreatedAt time.Time `json:"created_at"`
}