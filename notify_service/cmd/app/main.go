package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
	"traineesheep/notifyservice/internal/handlers"
	"traineesheep/notifyservice/internal/services"
	"traineesheep/notifyservice/internal/tgbot"
	"traineesheep/notifyservice/pkg/email"

	"github.com/tokyobordel/traineepkg/adapters/api/v1/auth"
	authMiddlewarePkg "github.com/tokyobordel/traineepkg/adapters/api/v1/middleware/authjwt"

	authSwagger "github.com/tokyobordel/traineepkg/adapters/api/v1/swagger"
	authService "github.com/tokyobordel/traineepkg/auth/service"
	jwtAuth "github.com/tokyobordel/traineepkg/authorization/jwt"

	"github.com/go-telegram/bot"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// SmtpDTO представляет собой структуру данных для конфигурации SMTP сервера.
// Используется для настройки параметров почтового сервера для отправки уведомлений.
type SmtpDTO struct {
	Email    string // Email адрес отправителя
	Password string // Пароль от email аккаунта
	Host     string // SMTP хост (например, smtp.gmail.com)
	Port     string // SMTP порт (например, 587)
}

// main является точкой входа в приложение.
//
// Функция выполняет следующие задачи:
//  1. Загружает переменные окружения из файла ./config/data.env
//  2. Устанавливает соединение с PostgreSQL базой данных
//  3. Инициализирует Telegram бота с обработчиком по умолчанию
//  4. Создает структуру SmtpDTO с параметрами почтового сервера
//  5. Подготавливает каналы и wait group для управления goroutines
//  6. Настраивает маршруты HTTP API с помощью Gorilla Mux
//  7. Запускает HTTP сервер и Telegram бота в отдельных goroutines
//  8. Ожидает сигнал завершения (SIGINT) для graceful shutdown
//  9. Корректно завершает работу всех компонентов
func main() {
	godotenv.Load("./config/data.env")
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	requiredEnvVars := []string{"DATABASE_CONNECT", "BOT_TOKEN", "WORKER_COUNT", "JWT_SECRET"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			logMsg := fmt.Sprintf("Required environment variable %s is not set", envVar)
			log.Fatal().Msg(logMsg)
		}
	}

	var channelSize, err = strconv.Atoi(os.Getenv("WORKER_COUNT"))
	if err != nil {
		log.Fatal().Msg("Error! Channel size can't be converted to int!")
		return
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_CONNECT"))
	if err != nil {
		log.Fatal().Msg("Error connecting to Database")
		return
	}
	defer pool.Close()

	opts := []bot.Option{
		bot.WithDefaultHandler(tgbot.Handler),
	}

	tgBot, err := bot.New(os.Getenv("BOT_TOKEN"), opts[0])
	if err != nil {
		log.Error().Msg("Error starting TG bot")
		return
	}
	smtpDto := email.NewSmtpDTO(os.Getenv("smtpEmail"), os.Getenv("smtpPassword"), os.Getenv("smtpHost"), os.Getenv("smtpPort"))

	grtChan := make(chan int, channelSize)
	wg := new(sync.WaitGroup{})
	d := handlers.NewDTO(tgBot, pool, smtpDto, grtChan, wg)

	app := fiber.New(fiber.Config{
		AppName: "Notify Service v1.0",
	})

	var myAuthService authService.IAuthService = services.NewUserService(os.Getenv("admin_password"), os.Getenv("admin_login"))
	jwtService := jwtAuth.NewService(os.Getenv("JWT_SECRET"), 15*time.Minute, 7*24*time.Hour)

	authMiddleware := authMiddlewarePkg.NewMiddleware(jwtService)
	handler := auth.NewHandler(myAuthService, jwtService, 15*time.Minute, 7*24*time.Hour)
	auth.SetupRouter(app, handler)
	authSwagger.SetupRouter(app)
	handlers.SetupNotifySwagger(app)
	api := app.Group("/api")

	api.Post("/notify", d.HandleNotify)

	// app.Get("/swagger/*", adaptor.HTTPHandler(httpSwagger.Handler(
	// 	httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	// )))

	protected := api.Group("/", authMiddleware.RequireAccessToken())

	protected.Post("/notify_types", d.HandleSaveSettingsCheckmarks)
	protected.Get("/get_notify_settings", d.HandleGetNotifySettings)
	serverErrors := make(chan error, 1)

	go func() {
		log.Info().Msg("Server has been started on port 8080!")
		serverErrors <- app.Listen(":8080")
	}()

	go func() {
		log.Info().Msg("Telegram bot has been started")
		tgBot.Start(context.Background())
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT)
	<-stop

	select {
	case err := <-serverErrors:
		logMsg := fmt.Sprintf("Server error: %v", err)
		log.Info().Msg(logMsg)
	case <-stop:
		log.Info().Msg("Shutting down server...")

		if err := app.Shutdown(); err != nil {
			logMsg := fmt.Sprintf("HTTP server shutdown error:", err)
			log.Error().Msg(logMsg)
		}

		log.Info().Msg("Waiting for active requests to finish...")
		wg.Wait()

		if _, err := tgBot.Close(context.Background()); err != nil {
			logMsg := fmt.Sprintf("Ошибка при отключении ТГ бота:", err)
			log.Error().Msg(logMsg)
		}
	}
}
