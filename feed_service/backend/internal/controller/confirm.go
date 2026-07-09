package controller

import (
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

// Confirm обрабатывает запрос на подтверждение аккаунта по токену из query-параметра "token".
// Токен парсится и проверяется через сервис токенов (TokenService),
// затем вызывается подтверждение учётной записи в UserService.
// При успешном подтверждении происходит редирект на корень "/".
func (ctrl *Controller) Confirm(c fiber.Ctx) error {
	// Получаем логгер, сохранённый в контексте запроса
	logger := c.Locals(utils.LoggerKey).(*zerolog.Logger)

	// Извлекаем токен из query-параметра "token"
	token := c.Query("token")

	// Парсим и валидируем токен подтверждения
	userID, userParseError := utils.ParseConfirmToken(token, ctrl.TokenService.GetSecret())
	if userParseError != nil {
		logger.Error().
			Str("token", token).
			Str("path", c.Path()).
			Msg("Ошибка парсинга токена: " + userParseError.Error())
		// Возвращаем ошибку в формате ApiResponse, но с Success: true (вероятно, требуется пересмотреть семантику)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Data:       "Некорректный токен",
			Success:    true,
			ErrMessage: "",
		})
	}

	// Подтверждаем учётную запись пользователя по его ID
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

	// При успехе перенаправляем пользователя на главную страницу
	return c.Redirect().To("/#login")
}
