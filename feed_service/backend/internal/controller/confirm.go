package controller

import (
	"log"
	"traineesheep/feedservice/internal/middleware"
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func (ctrl *Controller) Confirm(c *fiber.Ctx) error {
	token := c.Query("token")

	userID, userParseError := middleware.ParseToken(token)
	if userParseError != nil {
		log.Println(userParseError)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Data:       "Некорректный токен",
			Success:    true,
			ErrMessage: "",
		})
	}

	confirmErr := ctrl.UserService.ConfirmUserAccount(userID)
	if confirmErr != nil {
		log.Println(confirmErr)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Data:       "Некорректный токен",
			Success:    true,
			ErrMessage: "",
		})
	}

	log.Printf("GET /confirm: Аккаунт с user_id=%d подтвержден", userID)
	return c.Redirect("/")
}
