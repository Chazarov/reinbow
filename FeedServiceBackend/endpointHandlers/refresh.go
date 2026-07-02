package endpointHandlers

import (
	"traineesheep/feedservice/models"
	"traineesheep/feedservice/utils"
	"traineesheep/feedservice/jwtUtils"
	"github.com/gofiber/fiber/v2"
	"time"
)

func RefreshHandler(c *fiber.Ctx) error {
    // Парсим тело запроса
    var input struct {
        RefreshToken string `json:"refresh_token"`
    }
    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
            Success: false, ErrMessage: "Неверный формат запроса",
        })
    }
    if input.RefreshToken == "" {
        return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
            Success: false, ErrMessage: "refresh_token обязателен",
        })
    }

    // Ищем refresh-токен в БД
    var rt models.RefreshToken
    err := models.DB.QueryRow(
        "SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token = $1",
        input.RefreshToken,
    ).Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt)
    if err != nil {
        // Токен не найден (уже использован или недействителен)
        return c.Status(fiber.StatusUnauthorized).JSON(utils.ApiResponse{
            Success: false, ErrMessage: "Недействительный refresh-токен",
        })
    }

    // Проверяем срок действия
    if time.Now().After(rt.ExpiresAt) {
        // Удаляем просроченный токен
        models.DB.Exec("DELETE FROM refresh_tokens WHERE id = $1", rt.ID)
        return c.Status(fiber.StatusUnauthorized).JSON(utils.ApiResponse{
            Success: false, ErrMessage: "Refresh-токен истёк, выполните вход заново",
        })
    }

    // Токен валиден — удаляем его (ротация refresh-токенов)
    _, err = models.DB.Exec("DELETE FROM refresh_tokens WHERE id = $1", rt.ID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
            Success: false, ErrMessage: "Ошибка базы данных",
        })
    }

    // Получаем данные пользователя для ответа
    var user models.User
    err = models.DB.QueryRow(
        "SELECT id, username, created_at FROM users WHERE id = $1",
        rt.UserID,
    ).Scan(&user.ID, &user.Username, &user.CreatedAt)
    if err != nil {
        // Если пользователь вдруг не найден (маловероятно), вернём ошибку
        return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
            Success: false, ErrMessage: "Пользователь не найден",
        })
    }

    // Генерируем новую пару
    accessToken, err := jwtUtils.GenerateAccessToken(user)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
            Success: false, ErrMessage: "Не удалось создать access-токен",
        })
    }

    newRefreshToken, err := jwtUtils.GenerateRefreshToken()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
            Success: false, ErrMessage: "Не удалось создать refresh-токен",
        })
    }

    // Сохраняем новый refresh-токен в БД
    expiresAt := time.Now().Add(30 * 24 * time.Hour) // 30 дней
    _, err = models.DB.Exec(
        "INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
        rt.UserID, newRefreshToken, expiresAt,
    )
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
            Success: false, ErrMessage: "Не удалось сохранить refresh-токен",
        })
    }

    return c.JSON(utils.ApiResponse{
        Success: true,
        Data: fiber.Map{
            "access_token":  accessToken,
            "refresh_token": newRefreshToken,
        },
    })
}