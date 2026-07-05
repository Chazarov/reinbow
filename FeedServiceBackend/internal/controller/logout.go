package controller

import (
	"log"
	"time"
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// Logout обрабатывает POST-запрос на выход пользователя.
//
// Удаляет refresh-токен из базы данных (если он передан в куке `refresh_token`),
// затем очищает куки `refresh_token` и `access_token`, устанавливая их срок
// действия в прошлое. Куки устанавливаются как HttpOnly, SameSite=Strict.
// Ошибки при удалении токена из БД игнорируются — пользователь считается
// вышедшим в любом случае.
//
// Ответ всегда 200 OK:
//   - { success: true, data: null, err_message: "" }
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
