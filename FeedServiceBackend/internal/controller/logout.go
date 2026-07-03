package controller

import (
	"log"
	"time"
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func (ctrl *Controller) Logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	if refreshToken != "" {
		// Удаляем токен из БД (ошибки игнорируем — токен мог быть уже удалён)
		ctrl.TokenService.DeleteToken(refreshToken)
	}

	// Удаляем refresh_token и access_token в HttpOnly Secure куку
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/",
		Expires:  time.Unix(0, 0),
	})

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/",
		Expires:  time.Unix(0, 0),
	})

	log.Println("POST /logout: Пользователь разлогинен")
	return c.JSON(utils.ApiResponse{
		Success:    true,
		Data:       nil,
		ErrMessage: "",
	})
}
