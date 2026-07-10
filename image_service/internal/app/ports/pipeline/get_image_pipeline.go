// Пакет pipeline собирает цепочку обработки при получении изображения.
package pipeline

import (
	"context"
	"traineesheep/imageservice/internal/app/domain/interfaces"
	"traineesheep/imageservice/internal/app/models"
	"traineesheep/imageservice/internal/app/ports/connectors/imagickcon"
	"traineesheep/imageservice/internal/app/ports/connectors/rediscon"
	"traineesheep/imageservice/internal/app/ports/repository"

	"github.com/tokyobordel/traineepkg/errors"
	"github.com/tokyobordel/traineepkg/logger"
)

type imageGetterFunc func(
	ctx context.Context,
	mto models.GetImageMTO,
) ([]byte, errors.DomainError)

// GetImage реализует интерфейс IImageGetter для функции-обёртки.
func (f imageGetterFunc) GetImage(
	ctx context.Context,
	mto models.GetImageMTO,
) ([]byte, errors.DomainError) {
	return f(ctx, mto)
}

// NewGetImagePipeline собирает пайплайн получения изображения с обработкой и кешированием.
func NewGetImagePipeline(
	fileRepo *repository.ImageRepository,
	cacheConnector *rediscon.RedisCacheConnector,
	imagickConnector *imagickcon.ImagickConnector,
	log *logger.ContextLogger,
) interfaces.IImageGetter {
	getter := newDiskImageGetter(fileRepo)
	getter = withEditingToIcon(getter, imagickConnector, log)
	getter = withEditingBlur(getter, imagickConnector, log)
	getter = withEditingBlocked(getter, imagickConnector, log)
	getter = withCache(getter, cacheConnector, log)
	return getter
}

// newDiskImageGetter возвращает получатель, читающий изображение с диска.
func newDiskImageGetter(fileRepo *repository.ImageRepository) interfaces.IImageGetter {
	return imageGetterFunc(func(ctx context.Context, mto models.GetImageMTO) ([]byte, errors.DomainError) {
		if mto.Extension == nil {
			return nil, errors.NewInvalidParametersError("extension", nil, "extension is required")
		}
		return fileRepo.GetImage(ctx, mto.Id, *mto.Extension)
	})
}

// withEditingBlocked добавляет этап наложения баннера блокировки для гостевых запросов.
func withEditingBlocked(next interfaces.IImageGetter, imagickConnector *imagickcon.ImagickConnector, log *logger.ContextLogger) interfaces.IImageGetter {
	return imageGetterFunc(func(ctx context.Context, mto models.GetImageMTO) ([]byte, errors.DomainError) {
		data, derr := next.GetImage(ctx, mto)
		if derr != nil || imagickConnector == nil {
			return data, derr
		}
		if mto.MediaType == nil {
			return data, errors.NewInvalidParametersError("media_type", nil, "media type is required")
		}

		log.Infof(ctx, "IMG SIZE  EDIT BLOCKED: %.2f KB", float64(len(data))/1024)

		if mto.ImageStatus != nil && *mto.ImageStatus == models.Blocked && !mto.IsAdmin {
			data, derr = imagickConnector.EditingForTheBan(ctx, data, *mto.MediaType)
			if derr != nil {
				log.Criticalf(ctx, "Ошибка в пайплайне получения изображения. На этапе редактирования изображения(блюр): %v", derr.Error())
				return nil, derr
			}
			log.Infof(ctx, "Изображение было успешно отредактировано (наложение баннера блокировки) ")
		}

		return data, nil
	})
}

// withEditingBlur добавляет этап блюра и полосы модерации для гостевых запросов.
func withEditingBlur(next interfaces.IImageGetter, imagickConnector *imagickcon.ImagickConnector, log *logger.ContextLogger) interfaces.IImageGetter {
	return imageGetterFunc(func(ctx context.Context, mto models.GetImageMTO) ([]byte, errors.DomainError) {
		data, derr := next.GetImage(ctx, mto)
		if derr != nil || imagickConnector == nil {
			return data, derr
		}
		if mto.MediaType == nil {
			return data, errors.NewInvalidParametersError("media_type", nil, "media type is required")
		}

		log.Infof(ctx, "IMG SIZE  EDIT BLUR: %.2f KB", float64(len(data))/1024)

		if mto.ImageStatus != nil && *mto.ImageStatus == models.Unmoderated && !mto.IsAdmin {
			data, derr = imagickConnector.BlurWithModerationStripe(ctx, data, *mto.MediaType)
			if derr != nil {
				log.Criticalf(ctx, "Ошибка в пайплайне получения изображения. На этапе редактирования изображения(блюр): %v", derr.Error())
				return nil, derr
			}
			log.Infof(ctx, "Изображение было успешно отредактировано (блюр) ")
		}

		return data, nil
	})
}

// withEditingToIcon добавляет этап генерации иконки при запросе type=icon.
func withEditingToIcon(next interfaces.IImageGetter, imagickConnector *imagickcon.ImagickConnector, log *logger.ContextLogger) interfaces.IImageGetter {
	return imageGetterFunc(func(ctx context.Context, mto models.GetImageMTO) ([]byte, errors.DomainError) {
		data, derr := next.GetImage(ctx, mto)
		if derr != nil || imagickConnector == nil {
			return data, derr
		}
		if mto.MediaType == nil {
			return data, errors.NewInvalidParametersError("media_type", nil, "media type is required")
		}

		log.Infof(ctx, "IMG SIZE  EDIT ICON: %.2f KB", float64(len(data))/1024)

		if mto.ImageType == models.Icon {
			data, derr = imagickConnector.CompressToIcon(ctx, data, *mto.MediaType)
			if derr != nil {
				log.Criticalf(ctx, "Ошибка в пайплайне получения изображения. На этапе сжатия изображения: %v", derr)
				return nil, derr
			}
			log.Infof(ctx, "Изображение было успешно сжато ")
		}

		return data, nil
	})
}

// withCache оборачивает получатель слоем кеширования в Redis.
func withCache(
	next interfaces.IImageGetter,
	cacheConnector *rediscon.RedisCacheConnector,
	log *logger.ContextLogger,
) interfaces.IImageGetter {
	return imageGetterFunc(func(ctx context.Context, mto models.GetImageMTO) ([]byte, errors.DomainError) {

		return next.GetImage(ctx, mto)

		// if cacheConnector == nil {
		// 	return next.GetImage(ctx, mto)
		// }
		// imageCacheKey := rediscon.ImageCacheKey{
		// 	Id:        mto.Id,
		// 	IsAdmin:   mto.IsAdmin,
		// 	ImageType: mto.ImageType,
		// }
		// if cached, derr, hit := safeCacheGet(ctx, cacheConnector, imageCacheKey); derr != nil {
		// 	log.Errorf(ctx, "Cache read failed for image %d: %v", mto.Id, derr)
		// } else if hit {
		// 	log.Infof(ctx, "GETFROM CACHE Изображение %v получено из кеша", mto.Id)
		// 	return cached, nil
		// }
		// log.Infof(ctx, "NOTIN CACHE Изображение %v отсутствует в кеше", mto.Id)

		// data, derr := next.GetImage(ctx, mto)
		// if derr != nil {
		// 	return nil, derr
		// }

		// if derr := safeCacheSet(ctx, cacheConnector, imageCacheKey, data); derr != nil && log != nil {
		// 	log.Errorf(ctx, "Cache write failed for image %d: %v", mto.Id, derr)
		// }
		// log.Infof(ctx, "ADDTO CACHE Изображение %v отсутствует в кеше", mto.Id)

		// return data, nil
	})
}

// safeCacheGet безопасно читает изображение из кеша и возвращает признак попадания.
func safeCacheGet(ctx context.Context, cache *rediscon.RedisCacheConnector, imageCacheKey rediscon.ImageCacheKey) (data []byte, derr errors.DomainError, hit bool) {
	defer func() {
		if recover() != nil {
			data = nil
			derr = nil
			hit = false
		}
	}()

	data, derr = cache.GetImageFromCache(ctx, imageCacheKey)
	if derr != nil || len(data) == 0 {
		return data, derr, false
	}
	return data, nil, true
}

// safeCacheSet безопасно сохраняет изображение в кеш.
func safeCacheSet(ctx context.Context, cache *rediscon.RedisCacheConnector, imageCacheKey rediscon.ImageCacheKey, data []byte) (derr errors.DomainError) {
	defer func() {
		if recover() != nil {
			derr = nil
		}
	}()

	return cache.SetImage(ctx, imageCacheKey, data)
}
