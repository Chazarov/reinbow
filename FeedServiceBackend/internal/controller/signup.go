package controller

import (
	"log"
	"traineesheep/feedservice/internal/utils"
	"traineesheep/feedservice/internal/utils/notify"

	"github.com/gofiber/fiber/v2"
)

func (ctrl *Controller) Signup(c *fiber.Ctx) error {
	input, parseError := utils.ParseUserData(c, true)
	if parseError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: parseError.Error(),
		})
	}

	isUserExists, dbError := ctrl.UserService.ExistsByUsername(input.Username)

	if dbError != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Ошибка базы данных",
		})
	}

	if isUserExists {
		return c.Status(fiber.StatusConflict).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Пользователь с таким именем уже существует",
		})
	}

	user, userError := ctrl.UserService.CreateUser(input)

	if userError != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Не удалось создать пользователя. ",
		})
	}

	notify.NotifyUserRegistered(input.Username, input.Email, input.Password)

	log.Printf("POST /signup: Пользователь %s зарегистрирован", user.Username)
	// Успешная регистрация – возвращаем созданного пользователя
	return c.Status(fiber.StatusCreated).JSON(utils.ApiResponse{
		Data:       user,
		Success:    true,
		ErrMessage: "",
	})
}
