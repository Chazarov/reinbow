package service

import (
	"context"
	"traineesheep/imageservice/internal/app/domain/interfaces"
	"traineesheep/imageservice/internal/app/models"
	"traineesheep/imageservice/internal/app/ports/connectors/rediscon"
	"traineesheep/imageservice/internal/app/ports/repository"

	"github.com/tokyobordel/traineepkg/errors"
	"github.com/tokyobordel/traineepkg/logger"
)

// ImageService реализует бизнес-логику работы с изображениями.
type ImageService struct {
	relativeRepo *repository.ImagePostgresRepository
	fileRepo     *repository.ImageRepository
	cacheRepo    *rediscon.RedisCacheConnector
	imageGetter  interfaces.IImageGetter
	logger       *logger.ContextLogger
}

// NewImageService создаёт сервис работы с изображениями.
func NewImageService(
	relativeRepo *repository.ImagePostgresRepository,
	fileRepo *repository.ImageRepository,
	cacheRepo *rediscon.RedisCacheConnector,
	imageGetter interfaces.IImageGetter,
	logger *logger.ContextLogger,
) interfaces.IImageService {
	return &ImageService{
		relativeRepo: relativeRepo,
		fileRepo:     fileRepo,
		cacheRepo:    cacheRepo,
		imageGetter:  imageGetter,
		logger:       logger,
	}
}

// AddImage сохраняет метаданные и бинарное содержимое нового изображения.
func (s *ImageService) AddImage(ctx context.Context, createImageMTO models.CreateImageMTO) (models.Image, errors.DomainError) {
	if len(createImageMTO.Data) == 0 {
		s.logger.Errorf(ctx, "Image data cannot be empty")
		return models.Image{}, errors.NewInvalidParametersError("data", createImageMTO.Data, "image data cannot be empty")
	}

	if createImageMTO.Name == "" {
		s.logger.Errorf(ctx, "Image name cannot be empty")
		return models.Image{}, errors.NewInvalidParametersError("name", createImageMTO.Name, "image name cannot be empty")
	}

	if createImageMTO.MediaType == "" {
		s.logger.Errorf(ctx, "Media type cannot be empty")
		return models.Image{}, errors.NewInvalidParametersError("media_type", createImageMTO.MediaType, "media type cannot be empty")
	}

	image, err := s.relativeRepo.CreateImage(ctx, createImageMTO)
	if err != nil {
		return models.Image{}, err
	}

	extension, err := s.getExtensionByMediaType(createImageMTO.MediaType)
	if err != nil {
		s.logger.Errorf(ctx, "Unsupported media type: %s", createImageMTO.MediaType)
		return models.Image{}, errors.NewInvalidParametersError("media_type", createImageMTO.MediaType, "unsupported media type")
	}

	if err := s.fileRepo.PutImage(ctx, createImageMTO.Data, image.Id, extension); err != nil {
		s.logger.Errorf(ctx, "Failed to save image file for image %d: %v", image.Id, err)

		backErr := s.relativeRepo.DeleteImage(ctx, image.Id)
		if backErr != nil {
			s.logger.Criticalf(ctx, "Unable to roll back the image creation transaction after a file save failure image: %d: %v", image.Id, err)
		}
		return models.Image{}, err
	}

	s.logger.Infof(ctx, "Image added successfully with id: %d", image.Id)
	return image, nil
}

// GetImageContent возвращает бинарное содержимое изображения с учётом статуса и роли запроса.
func (s *ImageService) GetImageContent(ctx context.Context, mto models.GetImageMTO) ([]byte, string, errors.DomainError) {
	image, derr := s.relativeRepo.GetImage(ctx, mto.Id)
	if derr != nil {
		return nil, "", derr
	}

	extension, derr := s.getExtensionByMediaType(image.MediaType)
	if derr != nil {
		return nil, "", derr
	}

	mto.Extension = &extension
	mto.ImageStatus = &image.Status
	mediaType := models.MediaTypes(image.MediaType)
	mto.MediaType = &mediaType
	if mto.ImageType == "" {
		mto.ImageType = models.Usual
	}
	var data []byte
	data, derr = s.imageGetter.GetImage(
		ctx,
		mto,
	)
	if derr != nil {
		return nil, "", derr
	}

	return data, image.MediaType, nil
}

// GetImageMeta возвращает метаданные изображения без бинарного содержимого.
func (s *ImageService) GetImageMeta(ctx context.Context, id int) (models.Image, errors.DomainError) {

	image, derr := s.relativeRepo.GetImage(ctx, id)
	if derr != nil {
		return models.Image{}, derr
	}

	s.logger.Infof(ctx, "Image retrieved successfully: %d", image.Id)
	return image, nil
}

// GetImagesByStatus возвращает страницу изображений с указанным статусом модерации.
func (s *ImageService) GetImagesByStatus(ctx context.Context, status models.ModerStatus, pagination models.Pagination) ([]models.Image, errors.DomainError) {
	if status == "" {
		s.logger.Errorf(ctx, "Status cannot be empty")
		return nil, errors.NewInvalidParametersError("status", status, "status cannot be empty")
	}

	images, err := s.relativeRepo.GetImagesByStatus(ctx, status, pagination)
	if err != nil {
		return nil, err
	}

	s.logger.Infof(ctx, "Retrieved %d images with status %s", len(images), status)
	return images, nil
}

// GetAllImages возвращает страницу всех изображений.
func (s *ImageService) GetAllImages(ctx context.Context, pagination models.Pagination) ([]models.Image, errors.DomainError) {

	images, err := s.relativeRepo.GetAllImages(ctx, pagination)
	if err != nil {
		return nil, err
	}

	s.logger.Infof(ctx, "Retrieved %d images", len(images))
	return images, nil
}

// UpdateImageStatus обновляет статус модерации изображения и сбрасывает его кеш.
func (s *ImageService) UpdateImageStatus(ctx context.Context, id int, status models.ModerStatus) (models.Image, errors.DomainError) {
	image, derr := s.relativeRepo.UpdateImageStatus(ctx, id, status)
	if derr != nil {
		return models.Image{}, derr
	}

	if s.cacheRepo != nil {
		if cacheErr := s.cacheRepo.InvalidateImageCache(ctx, image.Id); cacheErr != nil {
			s.logger.Errorf(ctx, "Failed to invalidate image cache for image %d: %v", image.Id, cacheErr)
		}
	}

	return image, nil
}

// getExtensionByMediaType возвращает расширение файла по типу медиа.
func (s *ImageService) getExtensionByMediaType(mediaType string) (string, errors.DomainError) {
	mediaTypeMap := map[string]string{
		"jpeg": ".jpg",
		"png":  ".png",
		"gif":  ".gif",
		"webp": ".webp",
		"bmp":  ".bmp",
		"svg":  ".svg",
		"tiff": ".tiff",
		"ico":  ".ico",
		"heic": ".heic",
		"avif": ".avif",
	}

	ext, ok := mediaTypeMap[mediaType]
	if !ok {
		return "", errors.NewInvalidParametersError("media_type", mediaType, "unsupported media type")
	}

	return ext, nil
}

// GetImagesByStatusCount возвращает количество изображений с указанным статусом.
func (s *ImageService) GetImagesByStatusCount(ctx context.Context, status models.ModerStatus) (int, errors.DomainError) {
	return s.relativeRepo.GetImagesByStatusCount(ctx, status)
}

// GetImagesCount возвращает общее количество изображений в сервисе.
func (s *ImageService) GetImagesCount(ctx context.Context, status models.ModerStatus) (int, errors.DomainError) {
	return s.relativeRepo.GetImagesCount(ctx)
}
