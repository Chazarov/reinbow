package controller

import (
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
	"github.com/tokyobordel/traineepkg/adapters/api/v1/middleware/authjwt"
)

// NotifyAdmin обрабатывает запрос на отправку уведомления администратору во внешний сервис
// для пользователя, идентифицированного по JWT-токену.
// Идентификатор пользователя извлекается из контекста (ключ authjwt.UserIDContextKey).
// Вызывает метод SendNotificationExternal сервиса UserService.
// Возвращает JSON-ответ с сообщением об успехе или о том, что уведомление не отправлено.
func (ctrl *Controller) NotifyAdmin(c fiber.Ctx) error {
	// Получаем логгер, сохранённый в контексте запроса
	logger := c.Locals(utils.LoggerKey).(*zerolog.Logger)

	userID := c.Context().Value(authjwt.UserIDContextKey).(int)
	err := ctrl.UserService.SendNotificationExternal(userID)

	if err != nil {
		logger.Error().
			Str("path", c.Path()).
			Msg("Уведомление в сервис уведомлений не отправлено: " + err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Data:       "Уведомление в сервис уведомлений не отправлено",
			Success:    false,
			ErrMessage: "",
		})
	}

	logger.Error().
		Int("user_id", userID).
		Str("path", c.Path()).
		Msg("Уведомление в сервис уведомлений отправлено")
	return c.JSON(utils.ApiResponse{
		Data:       "Уведомление в сервис уведомлений отправлено",
		Success:    true,
		ErrMessage: "",
	})
}
