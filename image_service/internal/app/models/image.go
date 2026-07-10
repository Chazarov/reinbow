package models

import "time"

// Image представляет метаданные изображения, хранимые в базе данных.
type Image struct {
	Id        int         `json:"id" db:"id"`
	Name      string      `json:"name" db:"name"`
	MediaType string      `json:"media_type" db:"media_type"`
	CreatedAt time.Time   `json:"crated_at" db:"created_at"`
	Status    ModerStatus `json:"status" db:"status"`
}

// CreateImageMTO передаёт данные для создания нового изображения.
type CreateImageMTO struct {
	Data      []byte
	Name      string      `json:"name" db:"name"`
	MediaType string      `json:"media_type" db:"media_type"`
	Status    ModerStatus `json:"status" db:"status"`
}

// GetImageMTO передаёт параметры запроса бинарного содержимого изображения.
type GetImageMTO struct {
	Id          int       `json:"id"`
	ImageType   ImageType `json:"image_type"`
	Extension   *string
	MediaType   *MediaTypes
	ImageStatus *ModerStatus
	IsAdmin     bool
}
