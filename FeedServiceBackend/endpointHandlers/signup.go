package endpointHandlers

import (
	"fmt"
	"time"
	"traineesheep/feedservice/authUtils"
	"traineesheep/feedservice/models"
	"traineesheep/feedservice/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(c *fiber.Ctx) error {
	input, parseError := utils.ParseUserData(c, true)
	if parseError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: parseError.Error(),
		})
	}

	// Проверка существования пользователя
	exists := false
	err := models.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", input.Username).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Ошибка базы данных",
		})
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Пользователь с таким именем уже существует",
		})
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Не удалось обработать пароль",
		})
	}

	// Вставляем запись и получаем созданного пользователя
	var user models.User
	err = models.DB.QueryRow(
		"INSERT INTO users (username, password, created_at, tg_chat_id, email) VALUES ($1, $2, $3, $4, $5) RETURNING id, username, created_at, tg_chat_id, email",
		input.Username, string(hashedPassword), time.Now(), input.TgChatId, input.Email,
	).Scan(&user.ID, &user.Username, &user.CreatedAt, &user.TgChatId, &user.Email)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Не удалось создать пользователя: ",
		})
	}

	authUtils.NotifyUserRegistered(user)

	// Успешная регистрация – возвращаем созданного пользователя
	return c.Status(fiber.StatusCreated).JSON(utils.ApiResponse{
		Data:       user,
		Success:    true,
		ErrMessage: "",
	})
}