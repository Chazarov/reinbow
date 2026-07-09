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
	us *service.UserService,
	fs *service.FeedService,
	as authService.IAuthService, ts *jwt.Service) {

	// Инициализация контроллера с внедрением зависимостей
	ctrl := &Controller{
		AuthService:  as,
		UserService:  us,
		FeedService:  fs,
		TokenService: ts,
	}

	app.Use(middleware.RequestIDMiddleware())

	requireAccessToken := authMiddleware.NewMiddleware(ts).RequireAccessToken()

	// Публичная лента
	app.Get("/feed", ctrl.Feed)

	// Аутентификация
	app.Post("/auth/login", ctrl.Signin)  // получение токенов
	app.Post("/auth/logout", ctrl.Logout) // сброс refresh-токена

	// Пользователи
	app.Post("/users", ctrl.Signup) // создание пользователя (регистрация)

	// Подтверждение email (переход по ссылке из письма)
	// Middleware ConfirmRequired проверяет валидность токена из query-параметра
	app.Get("/users/confirm", middleware.ConfirmRequired(ts.GetSecret()), ctrl.Confirm)

	// Повторная отправка подтверждения (требуется аутентификация)
	//app.Post("/users/me/confirmation", middleware.TokenAuth, ctrl.SendConfirm)
	app.Post("/users/me/confirmation", requireAccessToken, ctrl.SendConfirm)

	// Профиль текущего пользователя (требуется аутентификация)
	app.Get("/users/me", requireAccessToken, ctrl.GetUser)

	// Посты
	// Создание нового поста (загрузка изображения и текста)
	app.Post("/posts", requireAccessToken, ctrl.Upload)

	// Healthcheck
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(utils.ApiResponse{
			Success:    true,
			Data:       nil,
			ErrMessage: "",
		})
	})
}
