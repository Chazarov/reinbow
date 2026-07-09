// Package app предоставляет конфигурацию и создание экземпляра Fiber-приложения.
//
// Здесь задаются глобальные ограничения (например, максимальный размер тела
// запроса) и настройки CORS для взаимодействия с фронтендом.
package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

// New создаёт и настраивает новый экземпляр Fiber-приложения.
//
// Устанавливает ограничение на размер тела запроса (10 МБ, что позволяет
// загружать до трёх изображений по 2 МБ каждое) и добавляет middleware CORS.
// Поддержка учётных данных (credentials) включена.
func New() *fiber.App {
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // ограничение - 10 мб (3 картинки по 2 мб)
	})

	// Настройка cors
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // теперь срез
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))

	return app
}
