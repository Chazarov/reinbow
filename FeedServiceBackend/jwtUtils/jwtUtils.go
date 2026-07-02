package jwtUtils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
	"traineesheep/feedservice/models"
	"traineesheep/feedservice/utils"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(utils.GetEnv("FEED_SERVICE_JWT_SECRET", 
    "Vj1WlmufcUengSqzIINyliPacXQXbSj0YqfTSYI3iWZ"))

var AuthRequired = jwtware.New(jwtware.Config{
    SigningKey: jwtware.SigningKey{Key: jwtSecret},
    ErrorHandler: func(c *fiber.Ctx, err error) error {
        return c.Status(fiber.StatusUnauthorized).JSON(utils.ApiResponse{
            Success: false, ErrMessage: err.Error(),
        })
    },
})

func GenerateAccessToken(user models.User) (string, error) {
    claims := jwt.MapClaims{
        "sub": strconv.Itoa(user.ID),
        "exp": time.Now().Add(15 * time.Minute).Unix(),
        "iat": time.Now().Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func ParseAccessToken(tokenStr string) (int, error) {
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

func GenerateRefreshToken() (string, error) {
    b := make([]byte, 32)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(b), nil
}