package controller

import (
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

// Signup обрабатывает POST-запрос на регистрацию нового пользователя.
//
// Ожидает JSON с полями username, email и password.
// Проверяет уникальность имени пользователя, создаёт
// запись в базе данных и отправляет уведомление о регистрации.
//
// Возможные ответы:
//   - 201: { success: true, data: созданный_пользователь }
//   - 400: { success: false, err_message: описание ошибки разбора }
//   - 409: { success: false, err_message: "Пользователь с таким именем уже существует" }
//   - 500: { success: false, err_message: "Ошибка базы данных" или "Не удалось создать пользователя. " }
func (ctrl *Controller) Signup(c fiber.Ctx) error {
	logger := c.Locals(utils.LoggerKey).(*zerolog.Logger)

	input, parseError := utils.ParseUserData(c, true)
	if parseError != nil {
		logger.Warn().
			Str("path", c.Path()).
			Msg("Ошибка парсинга данных")
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: parseError.Error(),
		})
	}

	isUserExists, dbError := ctrl.UserService.ExistsByUsername(input.Username)

	if dbError != nil {
		logger.Error().
			Str("username", input.Username).
			Str("path", c.Path()).
			Msg("Запрос выборки из БД вернул ошибку: " + dbError.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Ошибка базы данных",
		})
	}

	if isUserExists {
		logger.Warn().
			Str("username", input.Username).
			Str("path", c.Path()).
			Msg("Пользователь уже существует")
		return c.Status(fiber.StatusConflict).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Пользователь с таким именем уже существует",
		})
	}

	isUserExists, dbError = ctrl.UserService.ExistsByEmail(input.Email)

	if dbError != nil {
		logger.Error().
			Str("email", input.Email).
			Str("path", c.Path()).
			Msg("Запрос выборки из БД вернул ошибку: " + dbError.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Ошибка базы данных",
		})
	}

	if isUserExists {
		logger.Warn().
			Str("email", input.Email).
			Str("path", c.Path()).
			Msg("Пользователь уже существует")
		return c.Status(fiber.StatusConflict).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Пользователь с таким email уже существует",
		})
	}

	user, userError := ctrl.AuthService.Register(input.Password, input.Username,
		map[string]string{"email": input.Email})
	err := ctrl.UserService.NotifyClient.NotifyRegisterForAdmin(input.Username, input.Email)
	if err != nil {
		logger.Error().
			Str("username", input.Username).
			Str("path", c.Path()).
			Msg("Ошибка отправки уведомления: " + err.Error())
	} else {
		logger.Info().
			Str("username", input.Username).
			Str("path", c.Path()).
			Msg("Уведомление отправлено ")
	}
	if userError != nil {
		logger.Error().
			Str("username", input.Username).
			Str("path", c.Path()).
			Msg("Ошибка регистрации пользователя: " + userError.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Не удалось создать пользователя. ",
		})
	}

	pair, err := ctrl.TokenService.GenerateTokenPair(user.ID)
	if err != nil {
		logger.Error().
			Int("user_id", user.ID).
			Msg("Сервис генерации токенов вернул ошибку: " + err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Не удалось создать пользователя. ",
		})
	}
	utils.SetTokens(c, pair.AccessToken, pair.RefreshToken,
		ctrl.TokenService.GetAccessTTL(),
		ctrl.TokenService.GetRefreshTTL())

	logger.Info().
		Int("user_id", user.ID).
		Msg("Пользователь создан")
	// Успешная регистрация – возвращаем созданного пользователя
	return c.Status(fiber.StatusCreated).JSON(utils.ApiResponse{
		Data: fiber.Map{
			"user": user,
		},
		Success:    true,
		ErrMessage: "",
	})
}
