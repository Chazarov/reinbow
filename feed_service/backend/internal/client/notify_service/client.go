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

type NotifyClient struct {
	BaseURL    string
	SMTPClient *smtp.SmtpClient
	Client     *http.Client
}

func NewNotifyClient(baseURL string, SMTPClient *smtp.SmtpClient) *NotifyClient {
	return &NotifyClient{
		BaseURL:    baseURL,
		SMTPClient: SMTPClient,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (nsc *NotifyClient) NotifyRegisterForAdmin(username string, email string) error {
	payload := map[string]interface{}{
		"notify_type": "user_register",
		"email":       email,
		"message":     username + " зарегистрирован",
	}

	return nsc.SendPayload(payload)
}

// Методы для внутренней отправки уведомлений через внутренний сервис, а не внешний
func (nsc *NotifyClient) SendConfirmationEmail(username string,
	email string, conifrmToken string) {
	emailHost := utils.GetEnv("EMAIL_HOST", "http://localhost:3000")
	url := emailHost + "/api/users/confirm?token=" + conifrmToken
	msg := fmt.Sprintf("%s, перейдите по ссылке для подтверждения регистрации перейдите по ссылке: %s",
		username, url)
	nsc.SMTPClient.SendMessage([]string{email}, msg, "user_confirm")
}

func (nsc *NotifyClient) SendUserdataEmail(username string, passedEmail string, confirmToken string) {

	msg := fmt.Sprintf("Вы зарегистрированы в сервисе. Логин: %s",
		username)
	nsc.SMTPClient.SendMessage([]string{passedEmail}, msg, "user_register")

	nsc.SendConfirmationEmail(username, passedEmail, confirmToken)
}

func (nsc *NotifyClient) SendPayload(payload map[string]interface{}) error {
	if nsc.BaseURL == "" {
		return fmt.Errorf("Уведомления отключены")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Ошибка формирования уведомления: %v", err))
	}

	resp, err := nsc.Client.Post(nsc.BaseURL+"/notify", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Ошибка формирования уведомления: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Уведомление не доставлено, статус: %d", resp.StatusCode)
	} else {
		return nil
	}
}
