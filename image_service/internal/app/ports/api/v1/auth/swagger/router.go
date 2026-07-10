// Пакет swagger регистрирует Swagger UI для auth API.
package swagger

import (
	swaggo "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"

	_ "github.com/tokyobordel/traineepkg/adapters/api/v1/docs"
)

// SetupRouter подключает Swagger UI auth API к Fiber-приложению.
func SetupRouter(app *fiber.App) {
	app.Get("/auth/swagger/*", swaggo.HandlerDefault)
}
