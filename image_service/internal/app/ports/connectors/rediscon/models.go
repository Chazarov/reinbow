package rediscon

import "traineesheep/imageservice/internal/app/models"

// ImageCacheKey определяет составной ключ записи изображения в кеше.
type ImageCacheKey struct {
	Id        int
	IsAdmin   bool
	ImageType models.ImageType
}
