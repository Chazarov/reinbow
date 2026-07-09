package controller

import (
	"traineesheep/feedservice/internal/middleware"
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

func (ctrl *Controller) Confirm(c fiber.Ctx) error {
	logger := c.Locals(utils.LoggerKey).(*zerolog.Logger)

	token := c.Query("token")

	userID, userParseError := middleware.ParseConfirmToken(token, ctrl.TokenService.GetSecret())
	if userParseError != nil {
		logger.Error().
			Str("token", token).
			Str("path", c.Path()).
			Msg("Ошибка парсинга токена: " + userParseError.Error())
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Data:       "Некорректный токен",
			Success:    true,
			ErrMessage: "",
		})
	}

	confirmErr := ctrl.UserService.ConfirmUserAccount(userID)
	if confirmErr != nil {
		logger.Error().
			Str("token", token).
			Str("path", c.Path()).
			Msg("Ошибка подтверждения учетной записи: " + confirmErr.Error())
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Data:       "Ошибка подтверждения учетной записи",
			Success:    true,
			ErrMessage: "",
		})
	}

	logger.Info().
		Int("user_id", userID).
		Str("path", c.Path()).
		Msg("Аккаунт подтвержден")

	return c.Redirect().To("/")
}
