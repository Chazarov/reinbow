// Package app предоставляет конфигурацию и создание экземпляра Fiber-приложения.
//
// Здесь задаются глобальные ограничения (например, максимальный размер тела
// запроса) и настройки CORS для взаимодействия с фронтендом.
package app

import (
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// New создаёт и настраивает новый экземпляр Fiber-приложения.
//
// Устанавливает ограничение на размер тела запроса (10 МБ, что позволяет
// загружать до трёх изображений по 2 МБ каждое) и добавляет middleware CORS.
// Разрешённые источники CORS берутся из переменной окружения FRONTEND_ALLOWED_URLS.
// Поддержка учётных данных (credentials) включена.
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
