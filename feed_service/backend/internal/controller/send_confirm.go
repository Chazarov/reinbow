package controller

import (
	"traineesheep/feedservice/internal/middleware"
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
	"github.com/tokyobordel/traineepkg/adapters/api/v1/middleware/authjwt"
)

func (ctrl *Controller) SendConfirm(c fiber.Ctx) error {
	logger := c.Locals(utils.LoggerKey).(*zerolog.Logger)

	userID := c.Context().Value(authjwt.UserIDContextKey).(int)

	user, userError := ctrl.AuthService.GetMe(userID)
	if userError != nil {
		logger.Error().
			Int("user_id", userID).
			Str("path", c.Path()).
			Msg("Ошибка выборки данных из БД: " + userError.Error())
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Data:       "Некорректные данные",
			Success:    true,
			ErrMessage: "",
		})
	}
	email, ok := user.Data["email"]
	if !ok {
		logger.Error().
			Str("email", email).
			Str("path", c.Path()).
			Msg("Некорректный Email")
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Data:       "Ошибка отправки уведомления на почту. Попробуйте позже",
			Success:    true,
			ErrMessage: "",
		})
	}
	token, err := middleware.GenerateConfirmToken(userID)
	if err != nil {
		logger.Error().
			Str("email", email).
			Str("path", c.Path()).
			Msg("Ошибка создания т: " + err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Data:       "Ошибка отправки уведомления на почту. Попробуй позже",
			Success:    true,
			ErrMessage: "",
		})
	}
	ctrl.UserService.NotifyClient.SendConfirmationEmail(user.Login, email, token)

	logger.Error().
		Str("username", user.Login).
		Str("path", c.Path()).
		Msg("Уведомление о подтверждении отправлено")
	return c.JSON(utils.ApiResponse{
		Data:       "Уведомление о подтверждении отправлено",
		Success:    true,
		ErrMessage: "",
	})
}
