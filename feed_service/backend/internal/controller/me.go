package controller

import (
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
	"github.com/tokyobordel/traineepkg/adapters/api/v1/middleware/authjwt"
)

// GetUser обрабатывает GET-запрос данных текущего залогиненного пользователя.
func (ctrl *Controller) GetUser(c fiber.Ctx) error {
	logger := c.Locals(utils.LoggerKey).(*zerolog.Logger)
	userID := c.Context().Value(authjwt.UserIDContextKey).(int)
	user, userError := ctrl.AuthService.GetMe(userID)

	if userError != nil {
		logger.Warn().
			Int("user_id", userID).
			Str("path", c.Path()).
			Msg("Пользователь отсутствует в БД")
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Некорректные данные",
		})
	}

	logger.Info().
		Str("username", user.Login).
		Str("path", c.Path()).
		Msg("Пользователь есть в БД")
	return c.JSON(utils.ApiResponse{
		Data: fiber.Map{
			"user": user,
		},
		Success:    true,
		ErrMessage: "",
	})
}
