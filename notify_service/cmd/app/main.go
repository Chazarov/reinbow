package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"traineesheep/notifyservice/internal/handlers"
	"traineesheep/notifyservice/internal/services"
	"traineesheep/notifyservice/internal/tgbot"
	"traineesheep/notifyservice/internal/types"
	"traineesheep/notifyservice/pkg/email"

	"github.com/go-telegram/bot"
	"github.com/tokyobordel/traineepkg/adapters/api/v1/auth"
	authMiddlewarePkg "github.com/tokyobordel/traineepkg/adapters/api/v1/middleware/authjwt"
	"golang.org/x/sync/errgroup"

	authSwagger "github.com/tokyobordel/traineepkg/adapters/api/v1/swagger"
	authService "github.com/tokyobordel/traineepkg/auth/service"
	jwtAuth "github.com/tokyobordel/traineepkg/authorization/jwt"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var err error

// main является точкой входа в приложение.
//
// Функция выполняет следующие задачи:
//  1. Загружает переменные окружения из файла ./config/data.env
//  2. Устанавливает соединение с PostgreSQL базой данных
//  3. Инициализирует Telegram бота с обработчиком по умолчанию
//  4. Подготавливает каналы и wait group для управления goroutines
//  5. Настраивает маршруты HTTP API с помощью Fiber
//  6. Запускает HTTP сервер и Telegram бота в отдельных goroutines
//  7. Ожидает сигнал завершения (SIGINT) для graceful shutdown
//  8. Корректно завершает работу всех компонентов
func main() {
	godotenv.Load("./config/data.env")
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	tempVal, errconv := strconv.Atoi(os.Getenv("TG_ID"))
	ctx := types.Ctx
	if errconv != nil {
		log.Fatal().Msg("Telegram ID не может содержать что-то кроме цифр")
	}
	handlers.TelegramID = int64(tempVal)

	requiredEnvVars := []string{"DATABASE_CONNECT", "BOT_TOKEN", "WORKER_COUNT", "JWT_SECRET"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			logMsg := fmt.Sprintf("Required environment variable %s is not set", envVar)
			log.Fatal().Msg(logMsg)
		}
	}

	workersCount, err := strconv.Atoi(os.Getenv("WORKER_COUNT"))
	if err != nil {
		log.Fatal().Msg("Error! WORKER_COUNT can't be converted to int!")
		return
	}

	// Создаем errgroup с контекстом и устанавливаем лимит параллелизма
	taskGroup, _ := errgroup.WithContext(types.Ctx)
	taskGroup.SetLimit(workersCount)
	handlers.TaskGroup = taskGroup

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_CONNECT"))
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
	email.SmtpEmail = os.Getenv("smtpEmail")
	email.SmtpPass = os.Getenv("smtpPassword")
	email.SmtpHost = os.Getenv("smtpHost")
	email.SmtpPort = os.Getenv("smtpHost")

	handlers.TgBot = tgBot
	types.SqlConnection = pool

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

	api.Post("/notify", handlers.HandleNotify)

	protected := api.Group("/", authMiddleware.RequireAccessToken())

	protected.Post("/notify_types", handlers.HandleSaveSettingsCheckmarks)
	protected.Get("/notify_settings", handlers.HandleGetNotifySettings)
	serverErrors := make(chan error)

	go func() {
		log.Info().Msg("Server has been started on port 8080!")
		serverErrors <- app.Listen(":8080")
	}()

	go func() {
		log.Info().Msg("Telegram bot has been started")
		tgBot.Start(ctx)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		logMsg := fmt.Sprintf("Server error: %v", err)
		log.Info().Msg(logMsg)
	case <-stop:
		log.Info().Msg("Shutting down server...")

		if err := app.Shutdown(); err != nil {
			log.Error().Msgf("HTTP server shutdown error: %v", err)
		}

		log.Info().Msg("Waiting for active tasks to finish...")
		if err := handlers.TaskGroup.Wait(); err != nil {
			log.Error().Msgf("Task group error: %v", err)
		}

		if _, err := tgBot.Close(ctx); err != nil {
			log.Error().Msgf("Ошибка при отключении ТГ бота: %v", err)
		}
	}
}
