// Package controller реализует HTTP-обработчики (handlers) и регистрирует
// маршруты приложения.
//
// Контроллер связывает слои сервисов с входящими запросами, обеспечивает
// обработку бизнес-логики и возврат ответов.
package controller

import (
	"traineesheep/feedservice/internal/middleware"
	"traineesheep/feedservice/internal/service"

	"github.com/gofiber/fiber/v2"
)

// Controller хранит зависимости сервисов, необходимые для обработки запросов.
//
// Поля структуры не экспортируются, доступ к ним осуществляется через методы
// контроллера.
type Controller struct {
	UserService  *service.UserService
	FeedService  *service.FeedService
	TokenService *service.TokenService
}

// Create регистрирует все маршруты API в приложении Fiber и связывает их
// с соответствующими методами контроллера.
//
// Параметры:
//   - app: экземпляр *fiber.App, в котором регистрируются маршруты.
//   - us: сервис для работы с пользователями.
//   - fs: сервис для работы с лентой постов.
//   - ts: сервис для работы с токенами (refresh-токены).
//
// Регистрируемые маршруты:
//   - GET  /feed    – общая лента постов (открытый).
//   - GET  /refresh – обновление access-токена (требуется middleware.RefreshTokenRequired).
//   - POST /signin  – вход (открытый).
//   - POST /signup  – регистрация (открытый).
//   - POST /logout  – выход (открытый, удаление refresh-токена).
//   - POST /upload  – загрузка изображений и создание постов (требуется middleware.AuthRequired).
func Create(app *fiber.App, us *service.UserService,
	fs *service.FeedService, ts *service.TokenService) {
	// Инициализация контроллера. Через него реализовано прокидывание контекста
	ctrl := &Controller{
		UserService:  us,
		FeedService:  fs,
		TokenService: ts,
	}

	// Загрузка всей отсортированной ленты
	app.Get("/feed", ctrl.Feed)

	// Загрузка отсортированной ленты пользователя
	app.Get("/refresh", middleware.RefreshTokenRequired, ctrl.Refresh)

	// Вход
	app.Post("/signin", ctrl.Signin)

	// Регистрация
	app.Post("/signup", ctrl.Signup)

	// Удаление токена у пользователя
	app.Post("/logout", ctrl.Logout)

	// Загрузка изображений
	app.Post("/upload", middleware.AuthRequired, ctrl.Upload)
}
