package rediscon

import (
	"context"
	"fmt"
	"traineesheep/imageservice/internal/app/models"

	"github.com/tokyobordel/traineepkg/errors"
	"github.com/tokyobordel/traineepkg/logger"

	"github.com/go-redis/redis/v8"
)

// RedisCacheConnector кеширует обработанные изображения в Redis.
type RedisCacheConnector struct {
	logger      *logger.ContextLogger
	redisClient *redis.Client
}

// RedisChaheConnector создаёт коннектор кеша изображений в Redis.
func RedisChaheConnector(redisClient *redis.Client) *RedisCacheConnector {
	return &RedisCacheConnector{
		redisClient: redisClient,
	}
}

// GetImageFromCache возвращает изображение из кеша по составному ключу.
func (c *RedisCacheConnector) GetImageFromCache(ctx context.Context, imageCacheKey ImageCacheKey) ([]byte, errors.DomainError) {
	redisKey := c.getRedisKey(imageCacheKey)

	data, err := c.redisClient.Get(ctx, redisKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // не найдено
		}
		c.logger.Criticalf(ctx, "Ошибка получения изображения из кеша: %s", err.Error())
		derr := errors.NewInternalServiceError("ошибка кеширования", err)
		return nil, derr
	}

	return data, nil

}

// SetImage сохраняет изображение в кеш.
func (c *RedisCacheConnector) SetImage(ctx context.Context, imageCacheKey ImageCacheKey, data []byte) errors.DomainError {
	redisKey := c.getRedisKey(imageCacheKey)
	err := c.redisClient.Set(ctx, redisKey, data, 0).Err()
	if err != nil {
		c.logger.Criticalf(ctx, "Ошибка сохранения изображения в кеш: %s", err.Error())
		derr := errors.NewInternalServiceError("ошибка кеширования", err)
		return derr
	}

	return nil
}

// InvalidateImageCache удаляет все варианты кеша для указанного изображения.
func (c *RedisCacheConnector) InvalidateImageCache(ctx context.Context, imageID int) errors.DomainError {
	for _, isAdmin := range []bool{false, true} {
		for _, imageType := range []models.ImageType{models.Usual, models.Icon} {
			key := ImageCacheKey{Id: imageID, IsAdmin: isAdmin, ImageType: imageType}
			if err := c.deleteImageFromCache(ctx, key); err != nil {
				return err
			}
		}
	}
	return nil
}

// deleteImageFromCache удаляет одну запись изображения из кеша.
func (c *RedisCacheConnector) deleteImageFromCache(ctx context.Context, imageCacheKey ImageCacheKey) errors.DomainError {
	redisKey := c.getRedisKey(imageCacheKey)
	if err := c.redisClient.Del(ctx, redisKey).Err(); err != nil {
		if c.logger != nil {
			c.logger.Criticalf(ctx, "Ошибка удаления изображения из кеша: %s", err.Error())
		}
		return errors.NewInternalServiceError("ошибка удаления из кеша", err)
	}

	return nil
}

// getRedisKey формирует строковый ключ Redis для изображения.
func (s *RedisCacheConnector) getRedisKey(imageCacheKey ImageCacheKey) string {
	return fmt.Sprintf("%d:%t:%s", imageCacheKey.Id, imageCacheKey.IsAdmin, imageCacheKey.ImageType)
}
