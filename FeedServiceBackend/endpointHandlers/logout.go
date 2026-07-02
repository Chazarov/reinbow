package endpointHandlers

import (
	"traineesheep/feedservice/models"
	"traineesheep/feedservice/utils"
	"github.com/gofiber/fiber/v2"
)

func LogoutHandler(c *fiber.Ctx) error {
    var input struct {
        RefreshToken string `json:"refresh_token"`
    }
    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
            Success: false, ErrMessage: "Неверный формат запроса",
        })
    }

    if input.RefreshToken != "" {
        // Удаляем токен из БД (ошибки игнорируем — токен мог быть уже удалён)
        models.DB.Exec("DELETE FROM refresh_tokens WHERE token = $1", input.RefreshToken)
    }

    return c.JSON(utils.ApiResponse{
        Success: true,
        Data:    nil,
        ErrMessage: "",
    })
}