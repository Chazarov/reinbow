package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v3"
	fiberLogger "github.com/gofiber/fiber/v3/middleware/logger"
	fiberRecover "github.com/gofiber/fiber/v3/middleware/recover"
	_ "github.com/lib/pq"

	"traineesheep/imageservice/internal/app/domain/service"
	v1 "traineesheep/imageservice/internal/app/ports/api/v1"
	"traineesheep/imageservice/internal/app/ports/connectors/imagickcon"
	"traineesheep/imageservice/internal/app/ports/connectors/rediscon"
	"traineesheep/imageservice/internal/app/ports/connectors/traineenotify"
	"traineesheep/imageservice/internal/app/ports/pipeline"
	"traineesheep/imageservice/internal/app/ports/repository"
	"traineesheep/imageservice/internal/config"

	authRouter "github.com/tokyobordel/traineepkg/adapters/api/v1/auth"
	"github.com/tokyobordel/traineepkg/adapters/api/v1/middleware"
	authjwt "github.com/tokyobordel/traineepkg/adapters/api/v1/middleware/authjwt"
	jwtAuth "github.com/tokyobordel/traineepkg/authorization/jwt"
	"github.com/tokyobordel/traineepkg/logger"
)

const shutdownTimeout = 10 * time.Second

// main является точкой входа HTTP-сервиса изображений.
func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	appLogger, err := logger.NewContextLogger(cfg.DefaultLogsPath, cfg.CriticalLogsPath, cfg.LoggerDebug)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	db, err := sql.Open("postgres", cfg.Dsn())
	if err != nil {
		appLogger.Criticalf(ctx, "Failed to open database connection: %v", err)
		log.Fatal(err)
	}
	defer db.Close()

	for {
		pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		err := db.PingContext(pingCtx)
		cancel()

		if err != nil {
			appLogger.Errorf(ctx, "Failed to ping database, retrying in 3 seconds: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		appLogger.Info(ctx, "Database connected successfully")
		break
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr(),
	})
	defer redisClient.Close()

	for {
		pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		err := redisClient.Ping(pingCtx).Err()
		cancel()

		if err != nil {
			appLogger.Errorf(ctx, "Failed to ping redis, retrying in 2 seconds: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		appLogger.Info(ctx, "Redis connected successfully")
		break
	}

	cacheConnector := rediscon.RedisChaheConnector(redisClient)

	imagePostgresRepo := repository.NewImagePostgresRepository(db, appLogger)
	clientPostgresRepo := repository.NewClientPostgreRepository(db, appLogger)

	imageFileRepo, err := repository.NewImageRepository(ctx, cfg.ImageStoragePath, appLogger)
	if err != nil {
		appLogger.Criticalf(ctx, "Failed to create image repository: %v", err)
		log.Fatal(err)
	}

	authSvc := service.NewPkgAuthService(clientPostgresRepo, appLogger)
	jwtService := jwtAuth.NewService(cfg.JwtSecret, cfg.JwtAccessTTL, cfg.JwtRefreshTTL)

	imagickConnector := imagickcon.NewImagickConnector(appLogger)
	imageGetterPipeline := pipeline.NewGetImagePipeline(
		imageFileRepo,
		cacheConnector,
		imagickConnector,
		appLogger,
	)

	imageService := service.NewImageService(imagePostgresRepo, imageFileRepo, cacheConnector, imageGetterPipeline, appLogger)

	notifyService, err := traineenotify.NewNotificatorService(cfg.NotifyServiceURL, cfg.ExternalURL, cfg.TgAdminId, appLogger)
	if err != nil {
		appLogger.Criticalf(ctx, "Failed to create notification service: %v", err)
		log.Fatal(err)
	}

	spreadMiddleware := middleware.NewSpreadMiddleware(appLogger)
	authMiddleware := authjwt.NewMiddleware(jwtService)

	handler := v1.NewHandler(
		imageService,
		notifyService,
		appLogger,
		authMiddleware,
		cfg,
	)
	authHandler := authRouter.NewHandler(authSvc, jwtService, cfg.JwtAccessTTL, cfg.JwtRefreshTTL)

	app := fiber.New(fiber.Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	})
	app.Use(fiberLogger.New())
	app.Use(fiberRecover.New())
	app.Use(spreadMiddleware.AddSpreadInContext())
	authRouter.SetupRouter(app, authHandler)
	handler.SetupRoutes(app)

	appLogger.Infof(ctx, "Server starting on %s", cfg.ServerAddr())

	serverErrors := make(chan error, 1)
	go func() {
		if err := app.Listen(cfg.ServerAddr()); err != nil {
			serverErrors <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		appLogger.Criticalf(ctx, "Failed to start server: %v", err)
		log.Fatal(err)
	case sig := <-quit:
		appLogger.Infof(ctx, "Shutdown signal received: %s", sig)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		appLogger.Criticalf(ctx, "Server shutdown failed: %v", err)
		log.Fatal(err)
	}

	appLogger.Info(ctx, "Server stopped gracefully")
}
