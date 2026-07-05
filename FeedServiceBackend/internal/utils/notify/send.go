package notify

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"traineesheep/feedservice/internal/utils"
)

// sendPayload отправляет сформированный payload в виде JSON на URL,
// указанный в переменной окружения NOTIFY_URL (внешний сервис уведомлений).
//
// Если NOTIFY_URL не задан, уведомления отключаются и функция завершается
// без ошибки. При ошибках сериализации или сети они логируются, но не
// прерывают выполнение вызывающего кода. Таймаут запроса — 5 секунд.
//
// Параметры:
//   - payload: map с данными уведомления (тип, email, сообщение и т.д.).
func sendPayload(payload map[string]interface{}) {
	notifyURL := utils.GetEnv("NOTIFY_URL", "")
	if notifyURL == "" {
		log.Println("Уведомления отключены")
		return
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Ошибка формирования уведомления: %v", err)
		return
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(notifyURL, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Printf("Ошибка отправки уведомления: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Printf("Уведомление не доставлено, статус: %d", resp.StatusCode)
	} else {
		log.Printf("Уведомление отправлено")
	}
}
