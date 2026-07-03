package controller

import (
	"log"
	"traineesheep/feedservice/internal/middleware"
	"traineesheep/feedservice/internal/utils"
	"traineesheep/feedservice/internal/utils/notify"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

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
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/",
		MaxAge:   10 * 60,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/",
		MaxAge:   5 * 60,
	})

	notify.NotifyLogin(input.Username, input.Email)

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
