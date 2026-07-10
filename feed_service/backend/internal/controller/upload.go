package controller

import (
	"fmt"
	"strings"
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
	"github.com/tokyobordel/traineepkg/adapters/api/v1/middleware/authjwt"
)

// Upload обрабатывает POST-запрос на загрузку изображений и создание поста.
//
// Ожидает Content-Type "multipart/form-data". Допускается до трёх файлов
// в поле "images", каждый размером не более 2 МБ. Файлы отправляются во
// внешний сервис (ImageService), который сохраняет изображения и возвращает
// их идентификаторы. После успешной загрузки всех изображений создаётся
// новый пост с заголовком и описанием, полученными из полей формы.
//
// Заголовки безопасности: требует валидный access-токен в куке для
// идентификации пользователя.
func (ctrl *Controller) Upload(c fiber.Ctx) error {
	logger := c.Locals(utils.LoggerKey).(*zerolog.Logger)

	// Проверяем Content-Type
	contentType := string(c.Request().Header.ContentType())
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		logger.Error().
			Str("path", c.Path()).
			Msg("Неправильный формат")
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Success: false, ErrMessage: "Ожидается multipart/form-data",
		})
	}

	// Парсим форму
	form, err := c.MultipartForm()
	if err != nil {
		logger.Error().
			Str("path", c.Path()).
			Msg("Неправильный формат данных")
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Success: false, ErrMessage: "Неверный формат данных",
		})
	}

	files := form.File["images"] // поле, в котором фронт отправляет файлы
	if len(files) == 0 {
		logger.Error().
			Str("path", c.Path()).
			Msg("Нет изображений")
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Success: false, ErrMessage: "Не выбрано ни одного изображения",
		})
	}
	if len(files) > 3 {
		logger.Warn().
			Str("path", c.Path()).
			Msg("Превышен лимит изображений")
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Success: false, ErrMessage: "Максимальное количество изображений — 3",
		})
	}

	const maxFileSize = 2 << 20 // 2 МБ
	for _, file := range files {
		if file.Size > maxFileSize {
			msg := fmt.Sprintf("Файл '%s' превышает 2 МБ", file.Filename)
			logger.Error().
				Str("path", c.Path()).
				Msg(msg)
			return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
				Success:    false,
				ErrMessage: msg,
			})
		}
	}

	title := c.FormValue("title")
	description := c.FormValue("description")

	// Если поле не передано — можно установить значение по умолчанию
	if title == "" {
		title = "Без названия"
	}

	userID := c.Context().Value(authjwt.UserIDContextKey).(int)

	post, postError := ctrl.FeedService.CreatePost(userID, title, description, files)
	if postError != nil {
		logger.Error().
			Str("path", c.Path()).
			Msg("Ошибка создания поста: " + postError.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Success:    false,
			Data:       nil,
			ErrMessage: "Ошибка создания поста",
		})
	}

	msg := fmt.Sprintf("POST /upload: Пользователь %d отправил %d фото, создан пост с ID=%d",
		userID,
		len(post.Images),
		post.ID)
	logger.Info().
		Int("user_id", userID).
		Str("path", c.Path()).
		Msg(msg)

	return c.Status(fiber.StatusCreated).JSON(utils.ApiResponse{
		Success: true,
		Data: fiber.Map{
			"post":      post,
			"image_ids": post.Images,
		},
	})
}
