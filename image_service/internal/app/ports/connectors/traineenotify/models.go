package traineenotify

// NotifyType определяет тип отправляемого уведомления.
type NotifyType string

const (
	AdminNewImg    = "admin_newImg"
	UserImgVerdict = "user_imgVerdict"
)

// NotifyRequest описывает тело запроса к сервису уведомлений.
type NotifyRequest struct {
	Message    string     `json:"message"`
	NotifyType NotifyType `json:"notify_type"`
	TelegramId int        `json:"telegram_id"`
}

// NotifyResponse описывает ответ сервиса уведомлений.
type NotifyResponse struct {
	Data       interface{} `json:"data"`
	Success    bool        `json:"success"`
	ErrMessage string      `json:"err_message"`
}
