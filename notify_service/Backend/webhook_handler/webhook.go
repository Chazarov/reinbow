package webhook_handler

import (
	"bytes"
	"log"
	"net/http"
	"time"
)

// Функция SendWebhookMessage используется для отправки уведомлений на заданный URL
func SendWebhookMessage(url string, body_byte []byte) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body_byte))
	if err != nil {
		log.Println("Error while trying to create new request", err.Error())
		return
	}

	_, err = client.Do(req)
	if err != nil {
		log.Println("Error while trying to exec request:", err.Error())
		return
	}

	log.Println("Я отправил запрос на сайт", url)
}
