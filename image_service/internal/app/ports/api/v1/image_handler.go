package v1

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
	"traineesheep/imageservice/internal/app/models"

	"github.com/tokyobordel/traineepkg/adapters/api/v1/response"
	"github.com/tokyobordel/traineepkg/errors"

	"github.com/gofiber/fiber/v3"
)

// AddImageRequest описывает тело запроса на загрузку изображения.
type AddImageRequest struct {
	Name      string `json:"name" binding:"required" example:"test_image"`
	MediaType string `json:"media_type" binding:"required" example:"jpeg"`
	Data      []byte `json:"data" binding:"required" swaggertype:"string" format:"byte" example:"<base64>"`
}

// AddImage godoc
// @Summary      Загрузка изображения
// @Description  Сохраняет новое изображение со статусом unmoderated
// @Tags         images
// @Accept       json
// @Produce      json
// @Param        request body SwaggerAddImageRequest true "Данные изображения"
// @Success      201 {object} APIResponse{data=ImageMetaResponse}
// @Failure      400 {object} APIResponse
// @Failure      500 {object} APIResponse
// @Router       /image/upload [post]
func (h *Handler) AddImage(c fiber.Ctx) error {

	var req AddImageRequest
	if err := c.Bind().Body(&req); err != nil || req.Name == "" || req.MediaType == "" || len(req.Data) == 0 {
		h.logger.Errorf(c.Context(), "Invalid add image request: %v", err)
		response.MakeErrorResponse(c, log.Default(), errors.NewInvalidParametersError("body", nil, "invalid request body"))
		return nil
	}

	createImageMTO := models.CreateImageMTO{
		Name:      req.Name,
		MediaType: req.MediaType,
		Data:      req.Data,
		Status:    models.Unmoderated,
	}

	image, err := h.imageService.AddImage(c.Context(), createImageMTO)
	if err != nil {
		response.MakeErrorResponse(c, log.Default(), err)
		return nil
	}

	reqCtx := c.Context()
	imageID := image.Id

	go func() {
		ctx := context.WithoutCancel(reqCtx)
		defer func() {
			if r := recover(); r != nil {
				h.logger.Criticalf(ctx, "notify panic: %v\n%s", r, debug.Stack())
			}
		}()
		if err := h.notficationService.NotifyImageModeration(ctx, imageID); err != nil {
			h.logger.Criticalf(ctx, "Failed to send notification: %v", err)
		}
	}()

	response.MakeSuccessResponseWithStatus(c, http.StatusCreated, fiber.Map{
		"id":         image.Id,
		"name":       image.Name,
		"media_type": image.MediaType,
		"status":     image.Status,
		"created_at": image.CreatedAt,
	})
	return nil
}

// GetImageByIdGuest godoc
// @Summary      Получение изображения (гостевой доступ)
// @Description  Возвращает бинарное содержимое изображения для неавторизованных пользователей
// @Tags         images
// @Produce      application/octet-stream
// @Param        id path int true "ID изображения"
// @Param        type query string false "Тип изображения (icon)"
// @Success      200 {file} binary
// @Failure      400 {object} APIResponse
// @Failure      404 {object} APIResponse
// @Router       /guest/image/{id} [get]
func (h *Handler) GetImageByIdGuest(c fiber.Ctx) error {
	return h.getImageByRole(c, false, imageTypeFromQuery(c))
}

// GetImageByIdAdmin godoc
// @Summary      Получение изображения (администратор)
// @Description  Возвращает бинарное содержимое изображения без гостевой обработки
// @Tags         images
// @Produce      application/octet-stream
// @Security     AccessToken
// @Param        id path int true "ID изображения"
// @Param        type query string false "Тип изображения (icon)"
// @Success      200 {file} binary
// @Failure      400 {object} APIResponse
// @Failure      401 {object} APIResponse
// @Failure      404 {object} APIResponse
// @Router       /admin/image/{id} [get]
func (h *Handler) GetImageByIdAdmin(c fiber.Ctx) error {
	return h.getImageByRole(c, true, imageTypeFromQuery(c))
}

// GetImageMetaById godoc
// @Summary      Метаданные изображения
// @Description  Возвращает метаинформацию об изображении без бинарного содержимого
// @Tags         images
// @Produce      json
// @Param        id path int true "ID изображения"
// @Success      200 {object} APIResponse{data=ImageMetaResponse}
// @Failure      400 {object} APIResponse
// @Failure      404 {object} APIResponse
// @Router       /image/meta/{id} [get]
func (h *Handler) GetImageMetaById(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Errorf(c.Context(), "Invalid image id: %s", idStr)
		response.MakeErrorResponse(c, log.Default(), errors.NewInvalidParametersError("id", idStr, "invalid image id"))
		return nil
	}

	image, err := h.imageService.GetImageMeta(c.Context(), id)
	if err != nil {
		response.MakeErrorResponse(c, log.Default(), err)
		return nil
	}

	response.MakeSuccessResponse(c, fiber.Map{
		"id":         image.Id,
		"name":       image.Name,
		"media_type": image.MediaType,
		"status":     image.Status,
		"created_at": image.CreatedAt,
	})
	return nil
}

// imageContentType возвращает HTTP Content-Type для типа медиа.
func imageContentType(mediaType string) string {
	switch mediaType {
	case "jpeg":
		return "image/jpeg"
	case "svg":
		return "image/svg+xml"
	case "jpg", "png", "gif", "webp", "bmp", "tiff", "ico", "heic", "avif":
		return "image/" + mediaType
	default:
		return "application/octet-stream"
	}
}

// BlockImage godoc
// @Summary      Блокировка изображения
// @Description  Устанавливает изображению статус blocked
// @Tags         images
// @Produce      json
// @Security     AccessToken
// @Param        id path int true "ID изображения"
// @Success      200 {object} APIResponse{data=ImageStatusResponse}
// @Failure      400 {object} APIResponse
// @Failure      401 {object} APIResponse
// @Failure      404 {object} APIResponse
// @Router       /image/{id}/block [put]
func (h *Handler) BlockImage(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Errorf(c.Context(), "Invalid image id: %s", idStr)
		response.MakeErrorResponse(c, log.Default(), errors.NewInvalidParametersError("id", idStr, "invalid image id"))
		return nil
	}

	image, err := h.imageService.UpdateImageStatus(c.Context(), id, models.Blocked)
	if err != nil {
		response.MakeErrorResponse(c, log.Default(), err)
		return nil
	}

	response.MakeSuccessResponse(c, fiber.Map{
		"id":         image.Id,
		"name":       image.Name,
		"status":     image.Status,
		"updated_at": time.Now(),
	})
	return nil
}

// ApprovedImage godoc
// @Summary      Одобрение изображения
// @Description  Устанавливает изображению статус approved
// @Tags         images
// @Produce      json
// @Security     AccessToken
// @Param        id path int true "ID изображения"
// @Success      200 {object} APIResponse{data=ImageStatusResponse}
// @Failure      400 {object} APIResponse
// @Failure      401 {object} APIResponse
// @Failure      404 {object} APIResponse
// @Router       /image/{id}/approve [put]
func (h *Handler) ApprovedImage(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Errorf(c.Context(), "Invalid image id: %s", idStr)
		response.MakeErrorResponse(c, log.Default(), errors.NewInvalidParametersError("id", idStr, "invalid image id"))
		return nil
	}

	image, err := h.imageService.UpdateImageStatus(c.Context(), id, models.Approved)
	if err != nil {
		response.MakeErrorResponse(c, log.Default(), err)
		return nil
	}

	response.MakeSuccessResponse(c, fiber.Map{
		"id":         image.Id,
		"name":       image.Name,
		"status":     image.Status,
		"updated_at": time.Now(),
	})
	return nil
}

// GetUnmoderatedImages godoc
// @Summary      Список немодерированных изображений
// @Description  Возвращает страницу изображений со статусом unmoderated
// @Tags         images
// @Produce      json
// @Security     AccessToken
// @Param        page query int false "Номер страницы (с нуля)"
// @Param        page_size query int false "Размер страницы"
// @Success      200 {object} APIResponse{data=UnmoderatedImagesResponse}
// @Failure      400 {object} APIResponse
// @Failure      401 {object} APIResponse
// @Router       /image/unmoderated [get]
func (h *Handler) GetUnmoderatedImages(c fiber.Ctx) error {
	pagination, derr := ParsePagination(c, h.logger, h.paginationConfig)
	if derr != nil {
		response.MakeErrorResponse(c, log.Default(), derr)
		return nil
	}

	images, err := h.imageService.GetImagesByStatus(c.Context(), models.Unmoderated, pagination)
	if err != nil {
		response.MakeErrorResponse(c, log.Default(), err)
		return nil
	}

	count, err := h.imageService.GetImagesByStatusCount(c.Context(), models.Unmoderated)
	if err != nil {
		response.MakeErrorResponse(c, log.Default(), err)
		return nil
	}

	response.MakeSuccessResponse(c, fiber.Map{
		"images":      images,
		"count":       len(images),
		"total_count": count,
	})
	return nil
}

// GetAllImages godoc
// @Summary      Список всех изображений с пагинацией
// @Description  Возвращает все изображения в сервисе
// @Tags         images
// @Produce      json
// @Param        page query int false "Номер страницы (с нуля)"
// @Param        page_size query int false "Размер страницы"
// @Success      200 {object} APIResponse{data=UnmoderatedImagesResponse}
// @Failure      400 {object} APIResponse
// @Router       /image [get]
func (h *Handler) GetAllImages(c fiber.Ctx) error {
	pagination, derr := ParsePagination(c, h.logger, h.paginationConfig)
	if derr != nil {
		response.MakeErrorResponse(c, log.Default(), derr)
		return nil
	}

	images, err := h.imageService.GetAllImages(c.Context(), pagination)
	if err != nil {
		response.MakeErrorResponse(c, log.Default(), err)
		return nil
	}

	count, err := h.imageService.GetImagesCount(c.Context(), models.Unmoderated)
	if err != nil {
		response.MakeErrorResponse(c, log.Default(), err)
		return nil
	}

	response.MakeSuccessResponse(c, fiber.Map{
		"images":      images,
		"count":       len(images),
		"total_count": count,
	})
	return nil
}

// imageTypeFromQuery определяет тип запрашиваемого изображения из query-параметра type.
func imageTypeFromQuery(c fiber.Ctx) models.ImageType {
	if strings.EqualFold(c.Query("type"), string(models.Icon)) {
		return models.Icon
	}

	return models.Usual
}

// getImageByRole возвращает бинарное содержимое изображения с учётом роли запрашивающего.
func (h *Handler) getImageByRole(c fiber.Ctx, isAdmin bool, imageType models.ImageType) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Errorf(c.Context(), "Invalid image id: %s", idStr)
		response.MakeErrorResponse(c, log.Default(), errors.NewInvalidParametersError("id", idStr, "invalid image id"))
		return nil
	}

	content, mediaType, err := h.imageService.GetImageContent(c.Context(), models.GetImageMTO{
		Id:        id,
		ImageType: imageType,
		IsAdmin:   isAdmin,
	})
	if err != nil {
		response.MakeErrorResponse(c, log.Default(), err)
		return nil
	}

	c.Status(http.StatusOK)
	c.Set("Content-Type", imageContentType(mediaType))
	return c.Send(content)
}
