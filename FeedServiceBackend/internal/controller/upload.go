package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	jwt "traineesheep/feedservice/internal/middleware"
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func mimeTypeByExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return "application/octet-stream"
	}
	mt := mime.TypeByExtension(ext)
	if mt == "" {
		return "application/octet-stream"
	}
	// возвращаем основную часть, например "image/jpeg"
	if i := strings.Index(mt, ";"); i != -1 {
		return mt[:i]
	}
	return mt
}

func mediaTypeToShort(mediaType string) string {
	parts := strings.SplitN(mediaType, "/", 2)
	if len(parts) < 2 {
		return mediaType
	}
	subtype := parts[1]
	// убираем возможный суффикс, например "svg+xml" -> "svg"
	if idx := strings.Index(subtype, "+"); idx != -1 {
		subtype = subtype[:idx]
	}
	return subtype
}

func (ctrl *Controller) Upload(c *fiber.Ctx) error {
	// Проверяем Content-Type
	contentType := string(c.Request().Header.ContentType())
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Success: false, ErrMessage: "Ожидается multipart/form-data",
		})
	}

	// Парсим форму
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Success: false, ErrMessage: "Неверный формат данных",
		})
	}

	files := form.File["images"] // поле, в котором фронт отправляет файлы
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Success: false, ErrMessage: "Не выбрано ни одного изображения",
		})
	}
	if len(files) > 3 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Success: false, ErrMessage: "Максимальное количество изображений — 3",
		})
	}

	const maxFileSize = 2 << 20 // 2 МБ
	for _, file := range files {
		if file.Size > maxFileSize {
			return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
				Success:    false,
				ErrMessage: fmt.Sprintf("Файл '%s' превышает 2 МБ", file.Filename),
			})
		}
	}

	// URL внешнего сервиса
	imageAddURL := utils.GetEnv("IMAGE_ADD_URL", "")
	if imageAddURL == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Success: false, ErrMessage: "ImageService не настроен",
		})
	}

	// Структура ответа от внешнего сервиса (предполагаем {"id": 123})
	type externalImageResponse struct {
		ID int `json:"id"`
	}

	var imageIDs []int // сюда соберём полученные id

	// Отправляем каждый файл во внешний сервис
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
				Success: false, ErrMessage: "Не удалось прочитать файл",
			})
		}

		fileBytes, err := io.ReadAll(file)
		file.Close() // закрываем сразу, как прочитали
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
				Success: false, ErrMessage: "Ошибка чтения файла",
			})
		}

		// Определяем media_type: сначала пробуем из заголовка, затем по расширению
		mediaType := fileHeader.Header.Get("Content-Type")
		if mediaType == "" {
			mediaType = mimeTypeByExtension(fileHeader.Filename)
		}

		shortMediaType := mediaTypeToShort(mediaType)
		if shortMediaType == "" {
			shortMediaType = "octet-stream"
		}

		// Кодируем содержимое в base64
		encodedData := base64.StdEncoding.EncodeToString(fileBytes)

		// Формируем тело JSON
		requestBody := map[string]string{
			"name":       strings.TrimSuffix(fileHeader.Filename, filepath.Ext(fileHeader.Filename)),
			"media_type": shortMediaType,
			"data":       encodedData,
		}
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
				Success: false, ErrMessage: "Ошибка подготовки данных",
			})
		}

		// Отправляем POST-запрос с Content-Type: application/json
		httpReq, err := http.NewRequest("POST", imageAddURL, bytes.NewReader(jsonBody))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
				Success: false, ErrMessage: "Внутренняя ошибка",
			})
		}
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(httpReq)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
				Success: false, ErrMessage: "Не удалось сохранить изображение в ImageService",
			})
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			log.Printf("Ошибка внешнего сервиса: статус %d, тело %s", resp.StatusCode, string(bodyBytes))
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
				Success: false, ErrMessage: "ImageService вернул ошибку",
			})
		}

		var extResp utils.ApiResponse
		if err := json.NewDecoder(resp.Body).Decode(&extResp); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
				Success: false, ErrMessage: "Неверный ответ от ImageService",
			})
		}

		// 1. Приводим Data к map[string]interface{}
		dataMap, ok := extResp.Data.(map[string]interface{})
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
				Success: false, ErrMessage: "Неверный формат данных от ImageService",
			})
		}

		// 2. Извлекаем id
		idRaw, exists := dataMap["id"]
		if !exists {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
				Success: false, ErrMessage: "Поле id отсутствует в ответе ImageService",
			})
		}

		// 3. Пробуем привести к float64 (стандартный числовой тип из JSON)
		idFloat, ok := idRaw.(float64)
		if !ok {
			// Иногда id может быть строкой или другим типом — обработаем и это
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
				Success: false, ErrMessage: "Неверный тип id в ответе ImageService",
			})
		}

		// 4. Преобразуем в int
		imageID := int(idFloat)

		imageIDs = append(imageIDs, imageID)
	}

	// Все файлы успешно отправлены, создаём пост и записи в image_post
	refreshToken := c.Cookies("refresh_token")
	userID, userIDError := jwt.ParseToken(refreshToken)
	if userIDError != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Success: false, ErrMessage: "Некорректный токен",
		})
	}

	title := c.FormValue("title")
	description := c.FormValue("description")

	// Если поле не передано — можно установить значение по умолчанию
	if title == "" {
		title = "Без названия"
	}

	post, postError := ctrl.FeedService.CreatePost(userID, title, description, imageIDs)
	if postError != nil {
		c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Success:    false,
			Data:       nil,
			ErrMessage: "Ошибка создания поста",
		})
	}

	log.Printf("POST /upload: Пользователь %s отправил %d фото, создан пост с ID=%d",
		userID,
		len(imageIDs),
		post.ID)
	return c.Status(fiber.StatusCreated).JSON(utils.ApiResponse{
		Success: true,
		Data: fiber.Map{
			"post":      post,
			"image_ids": imageIDs,
		},
	})
}
