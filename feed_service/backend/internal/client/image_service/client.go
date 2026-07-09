package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"traineesheep/feedservice/internal/utils"
)

type ImageClient struct {
	BaseURL string
	Client  *http.Client
}

func NewImageClient(baseURL string) *ImageClient {
	return &ImageClient{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// mimeTypeByExtension возвращает MIME-тип файла на основе его расширения.
// Если расширение отсутствует или не распознано, возвращает "application/octet-stream".
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

// mediaTypeToShort преобразует полный MIME-тип в короткий формат (например, "image/jpeg" → "jpeg").
// Если тип не содержит '/', возвращает исходную строку.
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

func (isc *ImageClient) SaveFiles(files []*multipart.FileHeader) ([]int, error) {
	var imageIDs []int // сюда соберём полученные id
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("Не удалось прочитать файл")
		}

		fileBytes, err := io.ReadAll(file)
		file.Close() // закрываем сразу, как прочитали
		if err != nil {
			return nil, fmt.Errorf("Ошибка чтения файла")
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
			return nil, fmt.Errorf("Ошибка подготовки данных")
		}

		// Отправляем POST-запрос с Content-Type: application/json
		httpReq, err := http.NewRequest("POST", isc.BaseURL+"/upload", bytes.NewReader(jsonBody))
		if err != nil {
			return nil, fmt.Errorf("Внутренняя ошибка")
		}
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := isc.Client.Do(httpReq)
		if err != nil {
			return nil, fmt.Errorf("Не удалось сохранить изображение в ImageService")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf(fmt.Sprintf("Ошибка внешнего сервиса: статус %d, тело %s",
				resp.StatusCode, string(bodyBytes)))
		}

		var extResp utils.ApiResponse
		if err := json.NewDecoder(resp.Body).Decode(&extResp); err != nil {
			return nil, fmt.Errorf("Неверный ответ от ImageService")
		}

		// 1. Приводим Data к map[string]interface{}
		dataMap, ok := extResp.Data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Неверный формат данных от ImageService")
		}

		// 2. Извлекаем id
		idRaw, exists := dataMap["id"]
		if !exists {
			return nil, fmt.Errorf("Поле id отсутствует в ответе ImageService")
		}

		// 3. Пробуем привести к float64 (стандартный числовой тип из JSON)
		idFloat, ok := idRaw.(float64)
		if !ok {
			// Иногда id может быть строкой или другим типом — обработаем и это
			return nil, fmt.Errorf("Неверный тип id в ответе ImageService")
		}

		// 4. Преобразуем в int
		imageID := int(idFloat)

		imageIDs = append(imageIDs, imageID)
	}

	return imageIDs, nil
}
