package controller

import (
	"fmt"
	"log"
	"strconv"
	"traineesheep/feedservice/internal/middleware"
	models "traineesheep/feedservice/internal/model"
	"traineesheep/feedservice/internal/utils"
	"traineesheep/feedservice/pkg/email"

	"github.com/gofiber/fiber/v2"
)

func (ctrl *Controller) SendConfirm(c *fiber.Ctx) error {
	userID := c.Query("user_id")

	if userID == "" {
		log.Println("Передан некорректный user_id")
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Data:       "Некорректные данные",
			Success:    true,
			ErrMessage: "",
		})
	}

	userIDInt, userIDError := strconv.Atoi(userID)
	if userIDError != nil {
		log.Println("Передан некорректный user_id")
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Data:       "Некорректные данные",
			Success:    true,
			ErrMessage: "",
		})
	}

	user, userError := ctrl.UserService.GetByID(userIDInt)
	if userError != nil {
		log.Println("Передан некорректный user_id")
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Data:       "Некорректные данные",
			Success:    true,
			ErrMessage: "",
		})
	}

	sendError := SendConfirmationEmail(user)
	if sendError != nil {
		log.Printf("GET /send_confirm: Ошибка отправки уведомления на почту %s")
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Data:       "Ошибка отправки уведомления на почту. Попробуй позж",
			Success:    true,
			ErrMessage: "",
		})
	}

	log.Printf("GET /send_confirm: Отправлено уведомление на почту %s", user.Email)
	return c.JSON(utils.ApiResponse{
		Data:       "Уведомление о подтверждении отправлено",
		Success:    true,
		ErrMessage: "",
	})
}

func SendConfirmationEmail(user models.User) error {
	smtpData, smtpError := utils.GetSMTPData()
	if smtpError != nil {
		return smtpError
	}

	// Создаем токен
	conifrmToken, err := middleware.GenerateConfirmToken(user.ID)

	if err != nil {
		log.Printf("Не удалось сгенерировать токен для пользователя %s[id=%d]", user.Username)
	} else {
		emailHost := utils.GetEnv("EMAIL_HOST", "http://localhost:3000")
		url := emailHost + "/api/confirm?token=" + conifrmToken
		msg := fmt.Sprintf("%s, перейдите по ссылке для подтверждения регистрации перейдите по ссылке: %s",
			user.Username, url)
		smtp := email.NewSmtpDTO(smtpData["SMTP_EMAIL"],
			smtpData["SMTP_PASSWORD"],
			smtpData["SMTP_HOST"],
			smtpData["SMTP_PORT"])
		smtp.SendMessage([]string{user.Email}, msg, "user_confirm")
	}

	return err
}
