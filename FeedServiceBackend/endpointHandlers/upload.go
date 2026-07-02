package endpointHandlers

import (
	"traineesheep/feedservice/models"
	"traineesheep/feedservice/utils"
	"github.com/gofiber/fiber/v2"
	"time"
	"strings"
	"fmt"
	"net/http"
	"bytes"
	"io"
	"log"
	"encoding/json"
    "path/filepath"
    "mime"
    "encoding/base64"
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

func UploadHandler(c *fiber.Ctx) error {
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
                Success: false,
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
    /*type externalImageResponse struct {
        ID int `json:"id"`
    }*/

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
            // Но в вашем случае, судя по примеру, это число
            return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
                Success: false, ErrMessage: "Неверный тип id в ответе ImageService",
            })
        }

        // 4. Преобразуем в int (или сразу в нужный вам тип, например, int64)
        imageID := int(idFloat)

        imageIDs = append(imageIDs, imageID)
    }

    // Все файлы успешно отправлены, создаём пост и записи в image_post
    tx, err := models.DB.Begin()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
            Success: false, ErrMessage: "Ошибка базы данных",
        })
    }
    defer tx.Rollback()

    userID := c.FormValue("user_id")
    title := c.FormValue("title")
    description := c.FormValue("description")

    // Если поле не передано — можно установить значение по умолчанию
    if title == "" {
        title = "Без названия"
    }

    var post models.Post
    err = tx.QueryRow(
        "INSERT INTO post (user_id, title, description) VALUES ($1, $2, $3) RETURNING id, user_id, title, description, created_at",
        userID, title, description, // user_id = 0, пока нет авторизации
    ).Scan(&post.ID, &post.UserID, &post.Title, &post.Description, &post.CreatedAt)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
            Success: false, ErrMessage: "Не удалось создать пост",
        })
    }

    for _, imgID := range imageIDs {
        _, err = tx.Exec(
            "INSERT INTO image_post (post_id, image_id) VALUES ($1, $2)",
            post.ID, imgID,
        )
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
                Success: false, ErrMessage: "Ошибка привязки изображения к посту",
            })
        }
    }

    if err := tx.Commit(); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
            Success: false, ErrMessage: "Ошибка сохранения данных",
        })
    }

	// todo отправить уведомление

    return c.Status(fiber.StatusCreated).JSON(utils.ApiResponse{
        Success: true,
        Data: fiber.Map{
            "post":      post,
            "image_ids": imageIDs,
        },
    })
}