package types

// Структура ResponseData используется для формирования ответа на запрос
type ResponseData struct {
	Success       bool   `json:"success"`     // Удачный ли запрос
	Error_message string `json:"err_message"` // Сообщение об ошибке
	Token         string `json:"token"`       // Токен
}

// Структура Recipent содержит данные, которые будут использоваться для отправки
// уведомления пользователю
type Recipent struct {
	Email       string `json:"email"`       // Почта пользователя
	Notify_Type string `json:"notify_type"` // Тип уведомления
	Message     string `json:"message"`     // Сообщение
}

// Структура NotifyTypeMessenger сдержит данные, куда отправлять конкретный тип уведомлений
type NotifyTypeMessenger struct {
	NotifyType   string   `json:"notify_type"`        // Тип уведмоления
	Description  string   `json:"notify_description"` // Описание уведомления
	WantEmail    bool     `json:"want_email"`         // Отправляем ли на почту
	WantTelegram bool     `json:"want_telegram"`      // Отправляем ли в телеграм
	WebhookUrls  []string `json:"webhook_urls"`       // массив URL
}

// Структура NotifyTypeMessengerList нужна для хранения данных, полученных с сайта
// настройки уведомлений
type NotifyTypeMessengerList struct {
	// Success     bool                  `json:"Success"`     // Успех или нет
	Data []NotifyTypeMessenger `json:"data"` // Массив, содержащий данные куда отправлять определенные уведомления
}

// Структура LoginData нужна для хранения данных, передаваемых с сайта при авторизации
type LoginData struct {
	Login    string `json:"login"`    // Логин пользователя
	Password string `json:"password"` // Пароль пользователя
}
