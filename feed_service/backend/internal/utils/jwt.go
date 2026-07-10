// Package utils содержит общие утилиты и вспомогательные структуры,
// используемые в разных частях приложения: формат ответа API, работа
// с переменными окружения, парсинг и валидация входящих данных.
package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateConfirmToken создаёт подписанный JWT token
// для подтверждения регистрации указанного пользователя. Токен действителен 24 часа.
func GenerateConfirmToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"sub": strconv.Itoa(userID),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(GetEnv("FEED_SERVICE_JWT_SECRET", "")))
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
