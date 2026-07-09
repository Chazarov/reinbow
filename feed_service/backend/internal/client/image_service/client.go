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

// ImageClient — клиент для взаимодействия с удалённым сервисом изображений.
// Содержит базовый URL сервиса и HTTP-клиент с таймаутом.
type ImageClient struct {
	BaseURL string       // корневой URL сервиса (например, "http://images:8080")
	Client  *http.Client // HTTP-клиент для выполнения запросов
}

// NewImageClient создаёт новый экземпляр ImageClient с заданным базовым URL.
// Внутренний HTTP-клиент настраивается с таймаутом 10 секунд.
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
	// Убираем возможные параметры (например, "; charset=utf-8"),
	// оставляем только основную часть "image/jpeg"
	if i := strings.Index(mt, ";"); i != -1 {
		return mt[:i]
	}
	return mt
}

// mediaTypeToShort преобразует полный MIME-тип в короткий формат
// (например, "image/jpeg" → "jpeg", "image/svg+xml" → "svg").
// Если тип не содержит '/', возвращает исходную строку.
func mediaTypeToShort(mediaType string) string {
	parts := strings.SplitN(mediaType, "/", 2)
	if len(parts) < 2 {
		return mediaType
	}
	subtype := parts[1]
	// Удаляем возможный суффикс "+xml" и т.п.
	if idx := strings.Index(subtype, "+"); idx != -1 {
		subtype = subtype[:idx]
	}
	return subtype
}

// SaveFiles отправляет несколько файлов (multipart.FileHeader) в сервис изображений
// и возвращает список идентификаторов созданных изображений.
// На каждый файл делается отдельный POST-запрос на <BaseURL>/upload.
func (isc *ImageClient) SaveFiles(files []*multipart.FileHeader) ([]int, error) {
	var imageIDs []int // сюда будем собирать полученные id

	for _, fileHeader := range files {
		// Открываем содержимое загруженного файла
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("Не удалось прочитать файл")
		}

		// Читаем всё содержимое в память (для больших файлов может быть проблемой)
		fileBytes, err := io.ReadAll(file)
		file.Close() // закрываем сразу, как прочитали
		if err != nil {
			return nil, fmt.Errorf("Ошибка чтения файла")
		}

		// Определяем media_type: сначала пробуем из заголовка Content-Type multipart-части,
		// затем по расширению имени файла
		mediaType := fileHeader.Header.Get("Content-Type")
		if mediaType == "" {
			mediaType = mimeTypeByExtension(fileHeader.Filename)
		}

		// Приводим тип к короткому формату (например, "jpeg", "png")
		shortMediaType := mediaTypeToShort(mediaType)
		if shortMediaType == "" {
			shortMediaType = "octet-stream"
		}

		// Кодируем содержимое в base64 для передачи в JSON
		encodedData := base64.StdEncoding.EncodeToString(fileBytes)

		// Формируем тело запроса
		requestBody := map[string]string{
			"name":       strings.TrimSuffix(fileHeader.Filename, filepath.Ext(fileHeader.Filename)),
			"media_type": shortMediaType,
			"data":       encodedData,
		}
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("Ошибка подготовки данных")
		}

		// Создаём POST-запрос с JSON-телом
		httpReq, err := http.NewRequest("POST", isc.BaseURL+"/upload", bytes.NewReader(jsonBody))
		if err != nil {
			return nil, fmt.Errorf("Внутренняя ошибка")
		}
		httpReq.Header.Set("Content-Type", "application/json")

		// Выполняем запрос
		resp, err := isc.Client.Do(httpReq)
		if err != nil {
			return nil, fmt.Errorf("Не удалось сохранить изображение в ImageService")
		}
		defer resp.Body.Close()

		// Проверяем успешность ответа (200 или 201)
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("Ошибка внешнего сервиса: статус %d, тело %s",
				resp.StatusCode, string(bodyBytes))
		}

		// Разбираем JSON-ответ
		var extResp utils.ApiResponse
		if err := json.NewDecoder(resp.Body).Decode(&extResp); err != nil {
			return nil, fmt.Errorf("Неверный ответ от ImageService")
		}

		// Извлекаем ID изображения из поля Data ответа
		dataMap, ok := extResp.Data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Неверный формат данных от ImageService")
		}

		idRaw, exists := dataMap["id"]
		if !exists {
			return nil, fmt.Errorf("Поле id отсутствует в ответе ImageService")
		}

		idFloat, ok := idRaw.(float64) // JSON-числа парсятся как float64
		if !ok {
			return nil, fmt.Errorf("Неверный тип id в ответе ImageService")
		}

		imageID := int(idFloat)
		imageIDs = append(imageIDs, imageID)
	}

	return imageIDs, nil
}
