package v1

// APIResponse описывает стандартный формат ответа API.
type APIResponse struct {
	Data       interface{} `json:"data"`
	Success    bool        `json:"success" example:"true"`
	ErrMessage string      `json:"err_message" example:""`
	SpreadID   string      `json:"spread_id" example:"uuid-запроса"`
}

// AuthCredentialsRequest описывает учётные данные для входа или регистрации.
type AuthCredentialsRequest struct {
	Login string `json:"login" example:"testuser"`
	Pass  string `json:"pass" example:"123456"`
}

// UserData содержит дополнительные поля профиля пользователя.
type UserData struct {
	CreatedAt string `json:"created_at" example:"2026-07-03T09:00:00Z"`
}

// User описывает профиль пользователя в ответах auth API.
type User struct {
	ID    int      `json:"id" example:"1"`
	Login string   `json:"login" example:"testuser"`
	Data  UserData `json:"data"`
}

// AuthResponse возвращается при успешной аутентификации.
type AuthResponse struct {
	User User `json:"user"`
}

// LogoutResponse возвращается при успешном выходе из системы.
type LogoutResponse struct {
	Message string `json:"message" example:"logout successful"`
}

// ImageMetaResponse описывает метаданные изображения в API-ответе.
type ImageMetaResponse struct {
	ID        int    `json:"id" example:"1"`
	Name      string `json:"name" example:"test_image"`
	MediaType string `json:"media_type" example:"jpeg"`
	Status    string `json:"status" example:"unmoderated"`
	CreatedAt string `json:"created_at" example:"2026-07-03T09:00:00Z"`
}

// ImageStatusResponse описывает результат изменения статуса изображения.
type ImageStatusResponse struct {
	ID        int    `json:"id" example:"1"`
	Name      string `json:"name" example:"test_image"`
	Status    string `json:"status" example:"blocked"`
	UpdatedAt string `json:"updated_at" example:"2026-07-03T09:00:00Z"`
}

// ImageListItem описывает элемент списка изображений.
type ImageListItem struct {
	ID        int    `json:"id" example:"1"`
	Name      string `json:"name" example:"test_image"`
	MediaType string `json:"media_type" example:"jpeg"`
	CratedAt  string `json:"crated_at" example:"2026-07-03T09:00:00Z"`
	Status    string `json:"status" example:"unmoderated"`
}

// UnmoderatedImagesResponse описывает ответ со списком немодерированных изображений.
type UnmoderatedImagesResponse struct {
	Images     []ImageListItem `json:"images"`
	Count      int             `json:"count" example:"1"`
	TotalCount int             `json:"total_count" example:"1"`
}

// SwaggerAddImageRequest описывает тело запроса загрузки для Swagger.
type SwaggerAddImageRequest struct {
	Name      string `json:"name" example:"test_image"`
	MediaType string `json:"media_type" example:"jpeg"`
	Data      string `json:"data" example:"<base64>"`
}
