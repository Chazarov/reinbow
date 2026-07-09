// Package middleware предоставляет middleware-компоненты для Fiber-приложения,
// включая аутентификацию и авторизацию на основе JWT.
//
// Использует два токена:
//   - access_token (срок 5 минут) — для доступа к защищённым эндпоинтам,
//   - refresh_token (срок 10 минут) — для обновления access_token.
//
// Оба токена передаются в куках. Middleware AuthRequired проверяет access_token,
// RefreshTokenRequired — refresh_token. Секретный ключ подписи берётся из
// переменной окружения FEED_SERVICE_JWT_SECRET.
package middleware

import (
	"fmt"
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
)

func ConfirmRequired(secret []byte) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: secret,
		},
		Extractor: extractors.Extractor{
			Extract: func(c fiber.Ctx) (string, error) {
				value := c.Query("token")
				if value == "" {
					return "", fmt.Errorf("Нет токена")
				}
				return value, nil
			},
		},
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
				Success:    false,
				ErrMessage: err.Error(),
			})
		},
	})
}
