package endpointHandlers

import (
	"traineesheep/feedservice/models"
	"traineesheep/feedservice/utils"
	"traineesheep/feedservice/jwtUtils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func SigninHandler(c *fiber.Ctx) error {
	input, parseError := utils.ParseUserData(c, false)
	if parseError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: parseError.Error(),
		})
	}

	var user models.User
	err := models.DB.QueryRow(
		"SELECT id, username, password, created_at FROM users WHERE username = $1",
		input.Username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		// Пользователь не найден или неверный пароль
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Неверное имя пользователя или пароль",
		})
	}

	checkErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if checkErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Неверное имя пользователя или пароль",
		})
	}

	accessToken, err := jwtUtils.GenerateAccessToken(user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Ошибка создания access_token",
		})
	}

	refreshToken, err := jwtUtils.GenerateRefreshToken()
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Ошибка создания refresh_token",
		})
	}


	// Успех – возвращаем токен
	return c.JSON(utils.ApiResponse{
		Data: fiber.Map{
            "access_token": 	accessToken,
            "refresh_token": 	refreshToken,
			"user":				user,
        },
		Success:    true,
		ErrMessage: "",
	})
}