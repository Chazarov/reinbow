// Пакет v1 содержит HTTP-обработчики internal API.
package v1

import (
	"traineesheep/imageservice/internal/app/domain/interfaces"
	"traineesheep/imageservice/internal/config"

	authSwagger "traineesheep/imageservice/internal/app/ports/api/v1/auth/swagger"
	internalswagger "traineesheep/imageservice/internal/app/ports/api/v1/swagger"

	authMiddleware "github.com/tokyobordel/traineepkg/adapters/api/v1/middleware/authjwt"
	"github.com/tokyobordel/traineepkg/adapters/api/v1/response"
	"github.com/tokyobordel/traineepkg/logger"

	"github.com/gofiber/fiber/v3"
)

// Handler обрабатывает HTTP-запросы internal API.
type Handler struct {
	imageService       interfaces.IImageService
	notficationService interfaces.INotifyService
	logger             *logger.ContextLogger
	authMiddleware     *authMiddleware.Middleware
	paginationConfig   config.PaginationConfig
}

// NewHandler создаёт HTTP-обработчик internal API.
func NewHandler(
	imageService interfaces.IImageService,
	notficationService interfaces.INotifyService,
	logger *logger.ContextLogger,
	authMiddleware *authMiddleware.Middleware,
	paginationConfig config.PaginationConfig,
) *Handler {
	return &Handler{
		imageService:       imageService,
		notficationService: notficationService,
		logger:             logger,
		authMiddleware:     authMiddleware,
		paginationConfig:   paginationConfig,
	}
}

// SetupAuthSwagger регистрирует маршрут Swagger UI для auth API.
func (h *Handler) SetupAuthSwagger(app *fiber.App) {
	authSwagger.SetupRouter(app)
}

// SetupInternalSwagger регистрирует маршрут Swagger UI для internal API.
func (h *Handler) SetupInternalSwagger(app *fiber.App) {
	internalswagger.SetupRouter(app)
}

// SetupRoutes регистрирует маршруты internal API и Swagger.
func (h *Handler) SetupRoutes(app *fiber.App) {
	h.SetupAuthSwagger(app)
	h.SetupInternalSwagger(app)

	api := app.Group("/api")

	api.Get("/health", func(c fiber.Ctx) error {
		response.MakeSuccessResponse(c, nil)
		return nil
	})

	api.Get("/guest/image/:id", h.GetImageByIdGuest)
	api.Get("/image/meta/:id", h.GetImageMetaById)
	api.Get("/image", h.GetAllImages)
	api.Post("/image/upload", h.AddImage)

	protected := api.Group("/", h.authMiddleware.RequireAccessToken())
	protected.Get("/admin/image/:id", h.GetImageByIdAdmin)
	protected.Get("/image/unmoderated", h.GetUnmoderatedImages)
	protected.Put("/image/:id/block", h.BlockImage)
	protected.Put("/image/:id/approve", h.ApprovedImage)
}
