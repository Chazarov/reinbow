// Package utils содержит общие утилиты и вспомогательные структуры,
// используемые в разных частях приложения: формат ответа API, работа
// с переменными окружения, парсинг и валидация входящих данных.
package utils

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

func SetTokens(c fiber.Ctx, accessToken string, refreshToken string,
	accessTokenTTL time.Duration, refreshTokenTTL time.Duration) {
	if accessToken != "" {
		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Strict",
			Path:     "/",
			Expires:  time.Now().Add(accessTokenTTL),
		})
	}

	if refreshToken != "" {
		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Strict",
			Path:     "/",
			Expires:  time.Now().Add(refreshTokenTTL),
		})
	}
}
