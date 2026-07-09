// Package utils содержит общие утилиты и вспомогательные структуры,
// используемые в разных частях приложения: формат ответа API, работа
// с переменными окружения, парсинг и валидация входящих данных.
package utils

import (
	"fmt"
	"net/mail"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
)

// ApiResponse – единый формат ответа API.
// Используется всеми обработчиками для унификации возвращаемых данных.
type ApiResponse struct {
	Data       interface{} `json:"data"`
	Success    bool        `json:"success"`
	ErrMessage string      `json:"err_message"`
}

// UserData – структура для передачи данных пользователя между слоями
// (контроллер → сервис → DAO).
type UserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// GetEnv возвращает значение переменной окружения key, или fallback,
// если переменная не задана или пуста.
func GetEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// ParseUserData извлекает и валидирует данные пользователя из тела запроса.
// Если validateEmail == true, дополнительно проверяет, что email имеет
// корректный формат (через mail.ParseAddress).
// Возвращает заполненную структуру UserData или ошибку с описанием.
func ParseUserData(c fiber.Ctx, validateEmail bool) (UserData, error) {
	// Структура для парсинга данных пользователя
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := c.Bind().Body(&input); err != nil {
		return UserData{}, fmt.Errorf("неверный формат запроса: %w", err)
	}

	// Простая валидация
	if input.Username == "" || input.Password == "" {
		return UserData{}, fiber.ErrBadRequest
	}

	if validateEmail {
		_, emailErr := mail.ParseAddress(input.Email)
		if emailErr != nil {
			return UserData{}, fmt.Errorf("укажите адрес электронной почты в корректном формате")
		}
	}

	return input, nil
}

func GetSMTPData() (map[string]string, error) {
	smtpHost := GetEnv("SMTP_HOST", "")
	smtpPort := GetEnv("SMTP_PORT", "")
	smtpEmail := GetEnv("SMTP_EMAIL", "")
	smtpPassword := GetEnv("SMTP_PASSWORD", "")

	if smtpHost == "" || smtpPort == "" || smtpEmail == "" || smtpPassword == "" {
		return nil, fmt.Errorf("Уведомления отключены. Обратитесь к администратору сайта")
	}

	return map[string]string{
		"SMTP_HOST":     smtpHost,
		"SMTP_PORT":     smtpPort,
		"SMTP_EMAIL":    smtpEmail,
		"SMTP_PASSWORD": smtpPassword,
	}, nil
}

type UserProfile struct {
	Email       string    `json:"email"`
	IsConfirmed bool      `json:"is_confirmed"`
	CreatedAt   time.Time `json:"created_at"`
}
