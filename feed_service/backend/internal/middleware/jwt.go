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
	"strconv"
	"time"
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/golang-jwt/jwt/v5"
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

// GenerateConfirmToken создаёт подписанный JWT token
// для подтверждения регистрации указанного пользователя. Токен действителен 24 часа.
func GenerateConfirmToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"sub": strconv.Itoa(userID),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(utils.GetEnv("FEED_SERVICE_JWT_SECRET", "")))
}

func ParseConfirmToken(tokenStr string, secret []byte) (int, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil,
				fmt.Errorf("неподдерживаемый метод подписи")
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("неверные claims")
	}
	sub, _ := claims["sub"].(string)
	id, _ := strconv.Atoi(sub)
	return id, nil
}
