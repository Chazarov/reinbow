package webhook_handler

import (
	"bytes"
	"net/http"
	"time"
)

// Функция SendWebhookMessage используется для отправки уведомлений на заданный URL
func SendWebhookMessage(url string, bodyByte []byte) error {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyByte))
	if err != nil {
		return err
	}

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
