package controller

import (
	"log"
	"traineesheep/feedservice/internal/middleware"
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v2"
)

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
			ErrMessage: "Неверное имя пользователя или пароль",
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
