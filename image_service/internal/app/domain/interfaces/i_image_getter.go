package interfaces

import (
	"context"

	"traineesheep/imageservice/internal/app/models"

	"github.com/tokyobordel/traineepkg/errors"
)

// IImageGetter описывает получение бинарного содержимого изображения.
type IImageGetter interface {
	GetImage(ctx context.Context, mto models.GetImageMTO) ([]byte, errors.DomainError)
}
