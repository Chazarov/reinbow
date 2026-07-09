package controller

import (
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

// Signin обрабатывает POST-запрос на вход пользователя.
// Возможные ответы:
//   - 200: { success: true, data: { access_token, refresh_token, user } }
//   - 400: { success: false, err_message: описание ошибки разбора }
//   - 401: { success: false, err_message: "Неверное имя пользователя или пароль" }
//   - 500: { success: false, err_message: "Неверное имя пользователя или пароль" }
func (ctrl *Controller) Signin(c fiber.Ctx) error {
	logger := c.Locals(utils.LoggerKey).(*zerolog.Logger)
	input, parseError := utils.ParseUserData(c, false)
	if parseError != nil {
		logger.Warn().
			Msg("Ошибка парсинга данных")
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: parseError.Error(),
		})
	}

	user, checkErr := ctrl.AuthService.Login(input.Password, input.Username)
	if checkErr != nil {
		logger.Warn().
			Str("username", input.Username).
			Str("path", c.Path()).
			Msg("Неверное имя пользователя или пароль")
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Неверное имя пользователя или пароль",
		})
	}

	pair, err := ctrl.TokenService.GenerateTokenPair(user.ID)
	if err != nil {
		logger.Error().
			Str("username", input.Username).
			Str("path", c.Path()).
			Msg("Ошибка создания токена: " + err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Ошибка создания токенов",
		})
	}

	utils.SetTokens(c, pair.AccessToken, pair.RefreshToken,
		ctrl.TokenService.GetAccessTTL(),
		ctrl.TokenService.GetRefreshTTL())

	logger.Info().
		Str("username", input.Username).
		Str("path", c.Path()).
		Msg("Пользователь вошел в аккаунт")
	// Успех – возвращаем токен
	return c.JSON(utils.ApiResponse{
		Data: fiber.Map{
			"access_token":  pair.AccessToken,
			"refresh_token": pair.RefreshToken,
			"user":          user,
		},
		Success:    true,
		ErrMessage: "",
	})
}
