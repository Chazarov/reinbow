// FeedServiceBackend — серверная часть сервиса "ИН100ГРАММ".
// Обеспечивает REST API для регистрации, входа, загрузки постов и получения ленты.
// Использует фреймворк Fiber и PostgreSQL.
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"traineesheep/feedservice/internal/app"
	ISC "traineesheep/feedservice/internal/client/image_service"
	NSC "traineesheep/feedservice/internal/client/notify_service"
	"traineesheep/feedservice/internal/controller"
	"traineesheep/feedservice/internal/database"
	"traineesheep/feedservice/internal/repository"
	"traineesheep/feedservice/internal/service"
	"traineesheep/feedservice/internal/utils"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/tokyobordel/traineepkg/authorization/jwt"
	"github.com/tokyobordel/traineepkg/smtp"
)

// main — точка входа приложения.
// Загружает конфигурацию из .env-файла (по умолчанию development.env),
// подключается к базе данных, выполняет миграции, настраивает слои
// (DAO → Service) и запускает HTTP-сервер на адресе, указанном в
// переменной окружения BACKEND_HOST (по умолчанию :8080).
func main() {
	config := utils.GetEnv("DEPLOYMENT_CONFIG", "development")

	godotenv.Load(config + ".env")

	log.Printf("FeedServiceBackend: загружен %s конфиг\n", config)

	db := database.New() // соединяемся с БД

	database.Migrate(db) // инициализируем таблицы

	app := app.New() // создаём и конфигурируем fiber-приложение

	// Настройка сервисов, клиентов и DAO
	smtpData, smtpError := utils.GetSMTPData()
	if smtpError != nil {
		log.Fatal(smtpError)
	}
	SMTPClient := smtp.NewSmtpClient(smtpData["SMTP_EMAIL"],
		smtpData["SMTP_PASSWORD"],
		smtpData["SMTP_HOST"],
		smtpData["SMTP_PORT"])
	notifyClient := NSC.NewNotifyClient(utils.GetEnv("NOTIFICATION_SERVICE_URL", ""), SMTPClient)
	userDAO := repository.NewUserDAO(db)
	userService := service.NewUserService(userDAO, notifyClient)

	imageClient := ISC.NewImageClient(utils.GetEnv("IMAGE_SERVICE_URL", ""))
	feedDAO := repository.NewFeedDAO(db)
	feedService := service.NewFeedService(feedDAO, imageClient)

	JWTSecret := utils.GetEnv("FEED_SERVICE_JWT_SECRET", "")
	tokenService := jwt.NewService(JWTSecret, 5*time.Minute, 10*time.Minute)
	authService := service.NewAuthService(userDAO, notifyClient)

	// Задаем маршрутизацию
	controller.Create(app, userService, feedService, authService, tokenService)

	// Запускаем сервер в горутине
	addr := utils.GetEnv("BACKEND_HOST", ":8080")
	go func() {
		if err := app.Listen(addr); err != nil {
			// Логируем ошибку, но не паникуем, если это не ErrServerClosed
			log.Printf("Ошибка при запуске или остановке сервера: %v\n", err)
		}
	}()

	// Ждём сигнала завершения (Ctrl+C или docker stop)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("Получен сигнал %v, начинаю graceful shutdown...\n", sig)

	if err := app.Shutdown(); err != nil {
		log.Printf("Ошибка при остановке Fiber: %v\n", err)
	} else {
		log.Println("Сервер остановлен корректно")
	}
}
