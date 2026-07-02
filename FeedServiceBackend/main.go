package main

import (
	"log"
	"traineesheep/feedservice/endpointHandlers"
	"traineesheep/feedservice/jwtUtils"
	"traineesheep/feedservice/models"
	"traineesheep/feedservice/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	config := utils.GetEnv("DEPLOYMENT_CONFIG", "development")
	godotenv.Load(config)
	
	log.Printf("FeedServiceBackend: загружен %s конфиг\n", config)

	models.InitDB()

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // ограничение - 10 мб (3 картинки по 2 мб)
	})

	allowedUrls := utils.GetEnv("FRONTEND_ALLOWED_URLS", "http://localhost:3000")

	// Настройка cors
	app.Use(cors.New(cors.Config{
        AllowOrigins:     allowedUrls,
        AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
        AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
        AllowCredentials: false, // JWT передаётся в заголовке, куки не нужны
    }))

	// Загрузка всей отсортированной ленты
	app.Get("/loadMainFeed", endpointHandlers.LoadMainFeedHandler)

	// Загрузка отсортированной ленты пользователя
	app.Get("/loadUserFeed/:userID", endpointHandlers.LoadUserFeedHandler)

	// Вход
	app.Post("/signin", endpointHandlers.SigninHandler)

	// Регистрация
	app.Post("/signup", endpointHandlers.SignupHandler)

	// Обновление токена
	app.Post("/refresh", endpointHandlers.RefreshHandler)

	// Удаление токена у пользователя
	app.Post("/logout", endpointHandlers.LogoutHandler)

	// Загрузка изображений
	app.Post("/upload", jwtUtils.AuthRequired, endpointHandlers.UploadHandler)

	app.Listen(utils.GetEnv("BACKEND_HOST", ":8080"))
}
