package utils

import (
	"fmt"
	"os"
	"net/mail"
	"github.com/gofiber/fiber/v2"
)

// ApiResponse – единый формат ответа API
type ApiResponse struct {
	Data       interface{} `json:"data"`
	Success    bool        `json:"success"`
	ErrMessage string      `json:"err_message"`
}

// UserData – структура для парсинга входных данных
type UserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email 	 string `json:"email"`
}

// Функция для доступа к env'ам
func GetEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Функция для парсинга логина и пароля User
func ParseUserData(c *fiber.Ctx, validateEmail bool) (UserData, error) {
	// Структура для парсинга данных пользователя
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email string `json:"email"`
	}

	if err := c.BodyParser(&input); err != nil {
		return UserData{}, fmt.Errorf("неверный формат запроса: %w", err)
	}

	// Простая валидация
	if input.Username == "" || input.Password == "" {
		return UserData{}, fiber.ErrBadRequest
	}

	if validateEmail {
		_, emailErr := mail.ParseAddress(input.Email);
		if emailErr != nil {
			return UserData{}, fmt.Errorf("укажите адрес электронной почты в корректном формате")
		}
	}

	return input, nil
}
