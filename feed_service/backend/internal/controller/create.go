// Package controller реализует HTTP-обработчики (handlers) REST API
// и регистрирует маршруты приложения.
//
// Контроллер связывает слои сервисов с входящими запросами, реализует
// стандартные паттерны REST (ресурсы, HTTP-методы, статусы) и возвращает
// ответы в формате JSON.
package controller

import (
	"traineesheep/feedservice/internal/middleware"
	"traineesheep/feedservice/internal/service"
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v3"
	authPackage "github.com/tokyobordel/traineepkg/adapters/api/v1/auth"
	authMiddleware "github.com/tokyobordel/traineepkg/adapters/api/v1/middleware/authjwt"
	authService "github.com/tokyobordel/traineepkg/auth/service"
	"github.com/tokyobordel/traineepkg/authorization/jwt"
)

// Controller хранит зависимости сервисов, необходимые для обработки запросов.
//
// Поля структуры экспортируются, что позволяет при необходимости получать
// к ним доступ (например, в middleware), но основная работа с сервисами
// ведётся через методы контроллера.
type Controller struct {
	AuthService  authService.IAuthService // сервис аутентификации
	UserService  *service.UserService     // сервис управления пользователями
	FeedService  *service.FeedService     // сервис ленты постов
	TokenService *jwt.Service             // сервис работы с токенами
}

// Create регистрирует все маршруты REST API в приложении Fiber и связывает их
// с соответствующими методами контроллера.
//
// Параметры:
//   - app: экземпляр *fiber.App для регистрации маршрутов.
//   - us: сервис пользователей.
//   - fs: сервис ленты.
//   - ts: сервис токенов.
//   - as: внешний сервис аутентификации (пакет traineepkg).
//
// Ресурсы и маршруты:
//
//	Аутентификация:
//	  POST /auth/login  – получение токенов доступа (вход)
//	  POST /auth/logout – инвалидация refresh-токена (выход)
//
//	Пользователи:
//	  POST /users                      – регистрация нового пользователя
//	  GET  /users/me                   – получение профиля текущего пользователя (требуется токен)
//	  GET  /users/confirm?token=...    – подтверждение email по ссылке из письма
//	  POST /users/me/confirmation      – повторная отправка письма подтверждения (требуется токен)
//
//	Посты:
//	  POST /posts – создание поста (с загрузкой изображения) (требуется токен)
//
//	Лента:
//	  GET /feed – получение общей ленты постов (открытый)
//
//	Служебное:
//	  GET /health – проверка работоспособности сервиса
func Create(app *fiber.App,
	userService *service.UserService,
	feedService *service.FeedService,
	authService authService.IAuthService,
	tokenService *jwt.Service) {

	// Инициализация контроллера с внедрением зависимостей
	ctrl := &Controller{
		AuthService:  authService,
		UserService:  userService,
		FeedService:  feedService,
		TokenService: tokenService,
	}

	app.Use(middleware.RequestIDMiddleware())

	requireAccessToken := authMiddleware.NewMiddleware(tokenService).RequireAccessToken()

	// ---------- Публичные маршруты ----------
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(utils.ApiResponse{Success: true})
	})
	app.Get("/api/feed", ctrl.Feed) // общая лента

	// ---------- Группа /auth ----------
	authPackage.SetupRouter(app, authPackage.NewHandler(authService, tokenService,
		tokenService.GetAccessTTL(), tokenService.GetRefreshTTL()))

	// ---------- Группа /users ----------
	users := app.Group("/api/users")

	// Подтверждение email (свой middleware для проверки query-токена)
	users.Get("/confirm", middleware.ConfirmRequired(tokenService.GetSecret()), ctrl.Confirm)

	users.Post("/notify", requireAccessToken, ctrl.NotifyAdmin)
	users.Post("/confirmation", requireAccessToken, ctrl.SendConfirm) // POST /users/me/confirmation

	// ---------- Группа /posts (все требуют access-токен) ----------
	posts := app.Group("/api/posts", requireAccessToken)
	posts.Post("/", ctrl.Upload)
}
