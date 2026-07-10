// Пакет interfaces содержит доменные контракты сервиса изображений.
package interfaces

import (
	"context"
	"traineesheep/imageservice/internal/app/models"

	"github.com/tokyobordel/traineepkg/errors"
)

// IImageService описывает операции сервиса работы с изображениями.
type IImageService interface {
	AddImage(ctx context.Context, createImageMTO models.CreateImageMTO) (models.Image, errors.DomainError)
	GetImageMeta(ctx context.Context, id int) (models.Image, errors.DomainError)
	GetImageContent(ctx context.Context, mto models.GetImageMTO) ([]byte, string, errors.DomainError)
	GetImagesByStatus(ctx context.Context, status models.ModerStatus, pagination models.Pagination) ([]models.Image, errors.DomainError)
	GetAllImages(ctx context.Context, pagination models.Pagination) ([]models.Image, errors.DomainError)
	UpdateImageStatus(ctx context.Context, id int, status models.ModerStatus) (models.Image, errors.DomainError)
	GetImagesByStatusCount(ctx context.Context, status models.ModerStatus) (int, errors.DomainError)
	GetImagesCount(ctx context.Context, status models.ModerStatus) (int, errors.DomainError)
}
