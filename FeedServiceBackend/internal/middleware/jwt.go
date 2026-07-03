package middleware

import (
	"fmt"
	"strconv"
	"time"
	"traineesheep/feedservice/internal/utils"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(utils.GetEnv("FEED_SERVICE_JWT_SECRET",
	"Vj1WlmufcUengSqzIINyliPacXQXbSj0YqfTSYI3iWZ"))

var AuthRequired = jwtware.New(jwtware.Config{
	SigningKey:  jwtware.SigningKey{Key: jwtSecret},
	TokenLookup: "cookie:access_token",
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ApiResponse{
			Success: false, ErrMessage: err.Error(),
		})
	},
})

var RefreshTokenRequired = jwtware.New(jwtware.Config{
	SigningKey:  jwtware.SigningKey{Key: jwtSecret},
	TokenLookup: "cookie:refresh_token",
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ApiResponse{
			Success: false, ErrMessage: err.Error(),
		})
	},
})

func GenerateAccessToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"sub": strconv.Itoa(userID),
		"exp": time.Now().Add(5 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неподдерживаемый метод подписи")
		}
		return jwtSecret, nil
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

func GenerateRefreshToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"sub": strconv.Itoa(userID),
		"exp": time.Now().Add(10 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
