package controller

import (
	"log"
	"traineesheep/feedservice/internal/middleware"
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Signin обрабатывает POST-запрос на вход пользователя.
//
// Ожидает JSON с полями username и password. Проверяет существование
// пользователя в базе и корректность пароля (сравнение с bcrypt-хешем).
// При успехе генерирует access-токен (срок 5 минут) и refresh-токен
// (срок 10 минут), устанавливает их в HttpOnly Secure куки.
// Также отправляет уведомление о входе (notify.NotifyLogin).
//
// Возможные ответы:
//   - 200: { success: true, data: { access_token, refresh_token, user } }
//   - 400: { success: false, err_message: описание ошибки разбора }
//   - 401: { success: false, err_message: "Неверное имя пользователя или пароль" }
//   - 500: { success: false, err_message: "Неверное имя пользователя или пароль" }
//   - 401: { success: false, err_message: "Ошибка создания access_token" }
//   - 401: { success: false, err_message: "Ошибка создания refresh_token" }
func (ctrl *Controller) Signin(c *fiber.Ctx) error {
	input, parseError := utils.ParseUserData(c, false)
	if parseError != nil {
		log.Println(parseError)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: parseError.Error(),
		})
	}

	user, userError := ctrl.UserService.GetByUsername(input.Username)

	if userError != nil {
		log.Println(userError)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Неверное имя пользователя или пароль",
		})
	}

	checkErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if checkErr != nil {
		log.Println(checkErr)
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Неверное имя пользователя или пароль",
		})
	}

	accessToken, err := middleware.GenerateAccessToken(user.ID)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Ошибка создания access_token",
		})
	}

	refreshToken, err := middleware.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ApiResponse{
			Data:       nil,
			Success:    false,
			ErrMessage: "Ошибка создания refresh_token",
		})
	}

	// Устанавливаем refresh_token и access_token в HttpOnly Secure куку
	utils.SetTokens(c, accessToken, refreshToken)

	log.Printf("POST /signin: Пользователь %s вошел в аккаунт", user.Username)
	// Успех – возвращаем токен
	return c.JSON(utils.ApiResponse{
		Data: fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"user":          user,
		},
		Success:    true,
		ErrMessage: "",
	})
}
