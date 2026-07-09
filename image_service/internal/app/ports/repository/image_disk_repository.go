package repository

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tokyobordel/traineepkg/errors"
	"github.com/tokyobordel/traineepkg/logger"
)

// ImageRepository хранит бинарное содержимое изображений на диске.
type ImageRepository struct {
	DataDirectory string
	logger        *logger.ContextLogger
}

// NewImageRepository создаёт файловый репозиторий изображений и каталог хранения.
func NewImageRepository(ctx context.Context, dataDirectory string, logger *logger.ContextLogger) (*ImageRepository, errors.DomainError) {
	if err := os.MkdirAll(dataDirectory, 0755); err != nil {
		logger.Criticalf(ctx, "Failed to create directory: %v", err)
		derr := errors.NewInternalServiceError("failed to create directory", err)
		return nil, derr
	}
	return &ImageRepository{
		DataDirectory: dataDirectory,
		logger:        logger,
	}, nil
}

// PutImage сохраняет бинарное содержимое изображения на диск.
func (r *ImageRepository) PutImage(ctx context.Context, data []byte, imageId int, extension string) errors.DomainError {
	if len(data) == 0 {
		derr := errors.NewInvalidParametersError("data", data, "image data can't be empty")
		r.logger.Errorf(ctx, "Image data can't be empty")
		return derr
	}

	filename := fmt.Sprintf("%d%s", imageId, extension)
	filePath := filepath.Join(r.DataDirectory, filename)

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		derr := errors.NewInternalServiceError("failed to write image file", err)
		r.logger.Errorf(ctx, "Failed to write image file %d%s: %v", imageId, extension, err)
		return derr
	}

	return nil
}

// GetImage читает бинарное содержимое изображения с диска.
func (r *ImageRepository) GetImage(ctx context.Context, imageId int, extension string) ([]byte, errors.DomainError) {
	filename := fmt.Sprintf("%d%s", imageId, extension)
	filePath := filepath.Join(r.DataDirectory, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		derr := errors.NewNotFoundError("image", imageId)
		r.logger.Errorf(ctx, "Image %d%s not found", imageId, extension)
		return nil, derr
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		derr := errors.NewInternalServiceError("failed to read image file", err)
		r.logger.Criticalf(ctx, "Failed to read image %d%s: %v", imageId, extension, err)
		return nil, derr
	}

	return data, nil
}

// DeleteImage удаляет файл изображения с диска.
func (r *ImageRepository) DeleteImage(ctx context.Context, imageId int, extension string) errors.DomainError {
	filename := fmt.Sprintf("%d%s", imageId, extension)
	filePath := filepath.Join(r.DataDirectory, filename)

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			derr := errors.NewNotFoundError("image", imageId)
			r.logger.Errorf(ctx, "Image %d%s not found for deletion", imageId, extension)
			return derr
		}
		derr := errors.NewInternalServiceError("failed to delete image file", err)
		r.logger.Errorf(ctx, "Failed to delete image %d%s: %v", imageId, extension, err)
		return derr
	}

	r.logger.Infof(ctx, "Image %d%s deleted successfully", imageId, extension)
	return nil
}

// ImageExists проверяет наличие файла изображения на диске.
func (r *ImageRepository) ImageExists(ctx context.Context, imageId int, extension string) (bool, errors.DomainError) {
	filename := fmt.Sprintf("%d%s", imageId, extension)
	filePath := filepath.Join(r.DataDirectory, filename)

	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	derr := errors.NewInternalServiceError("failed to check image existence", err)
	r.logger.Errorf(ctx, "Failed to check existence of image %d%s: %v", imageId, extension, err)
	return false, derr
}

// GetAllImages возвращает имена всех файлов изображений в каталоге хранения.
func (r *ImageRepository) GetAllImages(ctx context.Context) ([]string, errors.DomainError) {
	files, err := os.ReadDir(r.DataDirectory)
	if err != nil {
		derr := errors.NewInternalServiceError("failed to read directory", err)
		r.logger.Errorf(ctx, "Failed to read directory %s: %v", r.DataDirectory, err)
		return nil, derr
	}

	var imageFiles []string
	for _, file := range files {
		if !file.IsDir() {
			imageFiles = append(imageFiles, file.Name())
		}
	}

	r.logger.Infof(ctx, "Found %d images in directory", len(imageFiles))
	return imageFiles, nil
}
