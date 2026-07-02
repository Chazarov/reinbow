package authUtils

import (
	"traineesheep/feedservice/models"
	"traineesheep/feedservice/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"bytes"
)

// notifyUserRegistered отправляет данные о новом пользователе на NOTIFY_URL
func NotifyUserRegistered(user models.User) {
    notifyURL := utils.GetEnv("NOTIFY_URL", "")
    if notifyURL == "" {
        return // уведомления отключены
    }

    // todo переписать payload (скорее всего, будет не такой)
    payload := map[string]interface{}{
		"notify_type": "user_register",
		"message": "Пользователь " + user.Username + " создан",
		"telegram_id": 924956695,
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
        log.Printf("Уведомление о регистрации пользователя %s отправлено", user.Username)
    }
}