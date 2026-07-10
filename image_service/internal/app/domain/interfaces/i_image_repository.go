package interfaces

import (
	"context"

	"github.com/tokyobordel/traineepkg/errors"
)

// IImageRepository описывает файловое хранилище изображений.
type IImageRepository interface {
	PutImage(ctx context.Context, data []byte, imageId int, extension string) errors.DomainError
	GetImage(ctx context.Context, imageId int, extension string) ([]byte, errors.DomainError)
	DeleteImage(ctx context.Context, imageId int, extension string) errors.DomainError
	ImageExists(ctx context.Context, imageId int, extension string) (bool, errors.DomainError)
	GetAllImages(ctx context.Context) ([]string, errors.DomainError)
}
