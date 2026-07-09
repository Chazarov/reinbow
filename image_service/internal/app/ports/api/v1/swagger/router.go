// Пакет internalswagger регистрирует Swagger UI для internal API.
package internalswagger

import (
	swaggo "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"

	_ "traineesheep/imageservice/api"
)

// SetupRouter подключает Swagger UI internal API к Fiber-приложению.
func SetupRouter(app *fiber.App) {
	app.Get("/internal/swagger/*", swaggo.New(swaggo.Config{
		InstanceName: "internal",
	}))
}
