package app

import (
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func New() *fiber.App {
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // ограничение - 10 мб (3 картинки по 2 мб)
	})

	allowedUrls := utils.GetEnv("FRONTEND_ALLOWED_URLS", "http://localhost:3000, http://10.64.11.142:3000")

	// Настройка cors
	app.Use(cors.New(cors.Config{
        AllowOrigins:     allowedUrls,
        AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
        AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
    }))

	return app
}