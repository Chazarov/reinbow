package v1

import (
	"strconv"
	"traineesheep/imageservice/internal/app/models"
	"traineesheep/imageservice/internal/config"

	"github.com/tokyobordel/traineepkg/errors"
	"github.com/tokyobordel/traineepkg/logger"

	"github.com/gofiber/fiber/v3"
)

// ParsePagination извлекает и валидирует параметры пагинации из query-параметров запроса.
func ParsePagination(c fiber.Ctx, logger *logger.ContextLogger, paginationConfig config.PaginationConfig) (models.Pagination, errors.DomainError) {
	ctx := c.Context()
	params := paginationConfig.GetPaginationParams()
	page := params.DefaultPage
	pageSize := params.DefaultPageSize

	if pageStr := c.Query("page"); pageStr != "" {
		parsed, err := strconv.Atoi(pageStr)
		if err != nil {
			logger.Errorf(ctx, "Invalid page parameter: %s", pageStr)
			return models.Pagination{}, errors.NewInvalidParametersError("page", pageStr, "page must be a non-negative integer")
		}
		if parsed < 0 {
			logger.Errorf(ctx, "Page parameter cannot be negative: %d", parsed)
			return models.Pagination{}, errors.NewInvalidParametersError("page", parsed, "page cannot be negative")
		}
		page = parsed
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		parsed, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			logger.Errorf(ctx, "Invalid page_size parameter: %s", pageSizeStr)
			return models.Pagination{}, errors.NewInvalidParametersError("page_size", pageSizeStr, "page_size must be a positive integer")
		}
		if parsed <= 0 {
			logger.Errorf(ctx, "Page_size parameter must be positive: %d", parsed)
			return models.Pagination{}, errors.NewInvalidParametersError("page_size", parsed, "page_size must be positive")
		}
		pageSize = parsed
	}

	return models.Pagination{
		Page:     page,
		PageSize: pageSize,
	}, nil
}
