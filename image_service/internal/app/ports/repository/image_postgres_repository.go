// Пакет repository реализует хранение данных изображений и клиентов.
package repository

import (
	"context"
	"database/sql"
	"time"
	"traineesheep/imageservice/internal/app/models"

	"github.com/tokyobordel/traineepkg/errors"
	"github.com/tokyobordel/traineepkg/logger"
)

// ImagePostgresRepository хранит метаданные изображений в PostgreSQL.
type ImagePostgresRepository struct {
	db     *sql.DB
	logger *logger.ContextLogger
}

// NewImagePostgresRepository создаёт репозиторий метаданных изображений в PostgreSQL.
func NewImagePostgresRepository(db *sql.DB, logger *logger.ContextLogger) *ImagePostgresRepository {
	return &ImagePostgresRepository{
		db:     db,
		logger: logger,
	}
}

// CreateImage сохраняет метаданные нового изображения в базе данных.
func (r *ImagePostgresRepository) CreateImage(ctx context.Context, image models.CreateImageMTO) (models.Image, errors.DomainError) {
	query := `
	INSERT INTO images (name, media_type, status, created_at)
	VALUES ($1, $2, $3, $4)
	RETURNING id, name, media_type, created_at, status`

	var createdImage models.Image
	err := r.db.QueryRowContext(
		ctx,
		query,
		image.Name,
		image.MediaType,
		image.Status,
		time.Now()).Scan(
		&createdImage.Id,
		&createdImage.Name,
		&createdImage.MediaType,
		&createdImage.CreatedAt,
		&createdImage.Status,
	)

	if err != nil {
		r.logger.Errorf(ctx, "Failed to create image: %v", err)
		derr := errors.NewInternalServiceError("failed to create image", err)
		return models.Image{}, derr

	}

	return createdImage, nil
}

// GetImage возвращает метаданные изображения по идентификатору.
func (r *ImagePostgresRepository) GetImage(ctx context.Context, imageId int) (models.Image, errors.DomainError) {
	query := `
		SELECT id, name, media_type, created_at, status 
		FROM images 
		WHERE id = $1
	`

	var image models.Image
	err := r.db.QueryRowContext(ctx, query, imageId).Scan(
		&image.Id,
		&image.Name,
		&image.MediaType,
		&image.CreatedAt,
		&image.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Errorf(ctx, "Image with id %d not found", imageId)
			derr := errors.NewNotFoundError("image", imageId)
			return models.Image{}, derr
		}
		r.logger.Errorf(ctx, "Failed to get image %d: %v", imageId, err)
		derr := errors.NewInternalServiceError("failed to get image", err)
		return models.Image{}, derr
	}

	return image, nil
}

// UpdateImageStatus обновляет статус модерации изображения.
func (r *ImagePostgresRepository) UpdateImageStatus(ctx context.Context, imageId int, status models.ModerStatus) (models.Image, errors.DomainError) {
	query := `
		UPDATE images 
		SET status = $1 
		WHERE id = $2 
		RETURNING id, name, media_type, created_at, status
	`

	var updatedImage models.Image
	err := r.db.QueryRowContext(ctx, query, status, imageId).Scan(
		&updatedImage.Id,
		&updatedImage.Name,
		&updatedImage.MediaType,
		&updatedImage.CreatedAt,
		&updatedImage.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Errorf(ctx, "Image with id %d not found for update", imageId)
			derr := errors.NewNotFoundError("image", imageId)
			return models.Image{}, derr
		}
		r.logger.Errorf(ctx, "Failed to update image %d status: %v", imageId, err)
		derr := errors.NewInternalServiceError("failed to update image status", err)
		return models.Image{}, derr
	}

	r.logger.Infof(ctx, "Image %d status updated to %s", imageId, status)
	return updatedImage, nil

}

// GetImagesByStatus возвращает страницу изображений с указанным статусом.
func (r *ImagePostgresRepository) GetImagesByStatus(ctx context.Context, status models.ModerStatus, pagination models.Pagination) ([]models.Image, errors.DomainError) {
	query := `
		SELECT id, name, media_type, created_at, status 
		FROM images 
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.db.QueryContext(ctx, query, status, pagination.PageSize, pagination.Page*pagination.PageSize)
	if err != nil {
		r.logger.Errorf(ctx, "Failed to get images by status %s: %v", status, err)
		return nil, errors.NewInternalServiceError("failed to get images", err)
	}
	defer rows.Close()

	var images []models.Image
	for rows.Next() {
		var image models.Image
		err := rows.Scan(
			&image.Id,
			&image.Name,
			&image.MediaType,
			&image.CreatedAt,
			&image.Status,
		)
		if err != nil {
			r.logger.Errorf(ctx, "Failed to scan image: %v", err)
			return nil, errors.NewInternalServiceError("failed to scan image", err)
		}
		images = append(images, image)
	}

	if err = rows.Err(); err != nil {
		r.logger.Errorf(ctx, "Rows iteration error: %v", err)
		return nil, errors.NewInternalServiceError("failed to iterate images", err)
	}

	r.logger.Infof(ctx, "Found %d images with status %s", len(images), status)
	return images, nil
}

// GetAllImages возвращает страницу всех изображений.
func (r *ImagePostgresRepository) GetAllImages(ctx context.Context, pagination models.Pagination) ([]models.Image, errors.DomainError) {
	query := `
		SELECT id, name, media_type, created_at, status 
		FROM images 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2;
	`

	rows, err := r.db.QueryContext(ctx, query, pagination.PageSize, pagination.Page*pagination.PageSize)
	if err != nil {
		r.logger.Errorf(ctx, "Failed to get images: %v", err)
		return nil, errors.NewInternalServiceError("failed to get images", err)
	}
	defer rows.Close()

	var images []models.Image
	for rows.Next() {
		var image models.Image
		err := rows.Scan(
			&image.Id,
			&image.Name,
			&image.MediaType,
			&image.CreatedAt,
			&image.Status,
		)
		if err != nil {
			r.logger.Errorf(ctx, "Failed to scan image: %v", err)
			return nil, errors.NewInternalServiceError("failed to scan image", err)
		}
		images = append(images, image)
	}

	if err = rows.Err(); err != nil {
		r.logger.Errorf(ctx, "Rows iteration error: %v", err)
		return nil, errors.NewInternalServiceError("failed to iterate images", err)
	}

	r.logger.Infof(ctx, "Found %d ", len(images))
	return images, nil
}

// GetImagesByStatusCount возвращает количество изображений с указанным статусом.
func (r *ImagePostgresRepository) GetImagesByStatusCount(ctx context.Context, status models.ModerStatus) (int, errors.DomainError) {
	query := `
		SELECT COUNT(*) as total_count
		FROM images
		WHERE status = $1;
	`

	var total_count int

	err := r.db.QueryRowContext(ctx, query, status).Scan(&total_count)

	if err != nil {
		r.logger.Errorf(ctx, "Fail to get images count %s", err.Error())
	}
	return total_count, nil
}

// GetImagesCount возвращает общее количество изображений.
func (r *ImagePostgresRepository) GetImagesCount(ctx context.Context) (int, errors.DomainError) {
	query := `
		SELECT COUNT(*) as total_count
		FROM images;
	`

	var total_count int

	err := r.db.QueryRowContext(ctx, query).Scan(&total_count)

	if err != nil {
		r.logger.Errorf(ctx, "Fail to get images count %s", err.Error())
	}
	return total_count, nil
}

// DeleteImage удаляет запись изображения из базы данных.
func (r *ImagePostgresRepository) DeleteImage(ctx context.Context, imageId int) errors.DomainError {
	query := `DELETE FROM images WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, imageId)
	if err != nil {
		r.logger.Errorf(ctx, "Failed to delete image %d: %v", imageId, err)
		return errors.NewInternalServiceError("failed to delete image", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Errorf(ctx, "Failed to get rows affected for image %d: %v", imageId, err)
		return errors.NewInternalServiceError("failed to delete image", err)
	}

	if rowsAffected == 0 {
		r.logger.Errorf(ctx, "Image with id %d not found for deletion", imageId)
		return errors.NewNotFoundError("image", imageId)
	}

	r.logger.Infof(ctx, "Image %d deleted successfully from database", imageId)
	return nil
}
