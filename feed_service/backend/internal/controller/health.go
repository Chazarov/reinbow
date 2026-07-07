package controller

import (
	"log"
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// Health
//
// # Нужен для compose
//
// Ответ всегда 200 OK:
//   - { success: true, data: null, err_message: "" }
func (ctrl *Controller) Health(c *fiber.Ctx) error {
	log.Println("POST /health: успех")
	return c.JSON(utils.ApiResponse{
		Success:    true,
		Data:       nil,
		ErrMessage: "",
	})
}
