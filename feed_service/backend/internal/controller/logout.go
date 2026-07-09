package controller

import (
	"time"
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

// Logout обрабатывает POST-запрос на выход пользователя.
//
// Очищает куки `refresh_token` и `access_token`, устанавливая их срок
// действия в прошлое. Куки устанавливаются как HttpOnly, SameSite=Strict.
// Ошибки при удалении токена из БД игнорируются — пользователь считается
// вышедшим в любом случае.
//
// Ответ всегда 200 OK:
//   - { success: true, data: null, err_message: "" }
func (ctrl *Controller) Logout(c fiber.Ctx) error {
	logger := c.Locals(utils.LoggerKey).(*zerolog.Logger)

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

	logger.Info().
		Str("path", c.Path()).
		Msg("Пользователь разлогинен")

	return c.JSON(utils.ApiResponse{
		Success:    true,
		Data:       nil,
		ErrMessage: "",
	})
}
