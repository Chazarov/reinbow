package controller

import (
	"log"
	"traineesheep/feedservice/internal/middleware"
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// Refresh обрабатывает GET-запрос на обновление access-токена.
//
// Требует наличия валидного refresh-токена в куке `refresh_token`
// (проверяется middleware.RefreshTokenRequired на уровне маршрута).
// При успехе генерирует новый access-токен, устанавливает его в куку
// `access_token` (HttpOnly, SameSite=Strict, срок действия 5 минут) и
// возвращает объект пользователя вместе с токенами.
//
// Возможные ответы:
//   - 200: { success: true, data: { access_token, refresh_token, user } }
//   - 401: { success: false, err_message: "Ошибка создания access_token" }
//   - 500: { success: false, err_message: "Некорректный токен" }
//   - 500: { success: false, err_message: "Пользователь не существует" }
func (ctrl *Controller) Refresh(c *fiber.Ctx) error {
	refreshTokenOld := c.Cookies("refresh_token")
	userID, userIDError := middleware.ParseToken(refreshTokenOld)
	if userIDError != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Success: false, ErrMessage: "Некорректный токен",
		})
	}

	accessToken, err := middleware.GenerateAccessToken(userID)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Ошибка создания access_token",
		})
	}

	// Устанавливаем access_token в HttpOnly Secure куку
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/",
		MaxAge:   5 * 60,
	})

	user, userError := ctrl.UserService.GetByID(userID)

	if userError != nil {
		log.Println(userError)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Пользователь не существует",
		})
	}

	log.Printf("GET /refresh: Обновлен access_token пользователя %s", user.Username)
	// Успех – возвращаем токен
	return c.JSON(utils.ApiResponse{
		Data: fiber.Map{
			"access_token":  accessToken,
			"refresh_token": c.Cookies("refresh_token"),
			"user":          user,
		},
		Success:    true,
		ErrMessage: "",
	})
}
