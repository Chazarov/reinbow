package models

import (
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email  	  string    `json:"email"`
	CreatedAt time.Time `json:"post_id"`
	Password  string    `json:"-"`
}