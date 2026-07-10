package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"traineesheep/feedservice/internal/utils"

	"github.com/tokyobordel/traineepkg/smtp"
)

// NotifyClient — клиент для отправки уведомлений.
// Поддерживает два канала: внешний HTTP-сервис уведомлений и SMTP-рассылку писем.
type NotifyClient struct {
	BaseURL    string           // базовый URL внешнего сервиса уведомлений (если пусто – HTTP-уведомления отключены)
	SMTPClient *smtp.SmtpClient // клиент для отправки email через SMTP
	Client     *http.Client     // HTTP-клиент для запросов к внешнему сервису
}

// NewNotifyClient создаёт новый экземпляр NotifyClient.
// baseURL — адрес внешнего сервиса уведомлений.
// SMTPClient — готовый клиент для отправки писем.
// HTTP-клиент инициализируется с таймаутом 10 секунд.
func NewNotifyClient(baseURL string, SMTPClient *smtp.SmtpClient) *NotifyClient {
	return &NotifyClient{
		BaseURL:    baseURL,
		SMTPClient: SMTPClient,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// NotifyRegisterForAdmin отправляет уведомление администратору о регистрации нового пользователя.
// Использует внешний HTTP-сервис (через SendPayload).
func (nsc *NotifyClient) NotifyRegisterForAdmin(username string, email string) error {
	payload := map[string]interface{}{
		"notify_type": "user_register",
		"email":       email,
		"message":     username + " зарегистрирован",
	}

	return nsc.SendPayload(payload)
}

func createConfirmationMessage(username string, token string) string {
	emailHost := utils.GetEnv("EMAIL_HOST", "http://localhost:3000")
	url := emailHost + "/api/users/confirm?token=" + token
	msg := fmt.Sprintf("%s, перейдите по ссылке для подтверждения регистрации перейдите по ссылке: %s",
		username, url)

	return msg
}

// SendConfirmationEmail отправляет письмо со ссылкой подтверждения регистрации.
// Использует SMTP-клиент.
// conifrmToken (опечатка в названии параметра, правильно "confirmToken") — токен подтверждения.
func (nsc *NotifyClient) SendConfirmationEmail(username string, email string, confirmToken string) {
	// Формируем URL подтверждения из переменной окружения EMAIL_HOST (или localhost по умолчанию)
	msg := createConfirmationMessage(username, confirmToken)
	// Отправляем письмо через SMTP с типом "user_confirm"
	nsc.SMTPClient.SendMessage([]string{email}, msg, "user_confirm")
}

// SendUserdataEmail отправляет пользователю два письма:
// 1. Уведомление о регистрации с логином.
// 2. Письмо со ссылкой подтверждения (вызывает SendConfirmationEmail).
func (nsc *NotifyClient) SendUserdataEmail(username string, passedEmail string, confirmToken string) {

	// Первое письмо: данные для входа
	confMsg := createConfirmationMessage(username, confirmToken)
	msg := fmt.Sprintf("Вы зарегистрированы в сервисе. Логин: %s.\n%s",
		username, confMsg)
	nsc.SMTPClient.SendMessage([]string{passedEmail}, msg, "user_register")
}

// SendPayload отправляет произвольное уведомление через внешний HTTP-сервис.
// Если BaseURL не задан, возвращает ошибку «Уведомления отключены».
// В случае ошибки на любом этапе (маршалинг JSON, HTTP-запрос, статус ответа >= 400)
// возвращает соответствующую ошибку.
func (nsc *NotifyClient) SendPayload(payload map[string]interface{}) error {
	if nsc.BaseURL == "" {
		return fmt.Errorf("Уведомления отключены")
	}

	// Преобразуем payload в JSON
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Ошибка формирования уведомления: %v", err))
	}

	// Отправляем POST-запрос на endpoint /notify
	resp, err := nsc.Client.Post(nsc.BaseURL+"/notify", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Ошибка формирования уведомления: %v", err))
	}
	defer resp.Body.Close()

	// Если статус ответа >= 400, считаем доставку неудачной
	if resp.StatusCode >= 400 {
		return fmt.Errorf("Уведомление не доставлено, статус: %d", resp.StatusCode)
	}
	return nil
}
