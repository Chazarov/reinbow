package handlers

import (
	_ "traineesheep/notifyservice/docs"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
)

func SetupNotifySwagger(app *fiber.App) {
	app.Get("/internal/swagger/*", swaggo.New(swaggo.Config{
		InstanceName: "internal",
	}))
}
