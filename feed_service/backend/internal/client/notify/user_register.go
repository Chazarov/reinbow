package notify

// NotifyRegister - отправляет уведомление о регистрации пользователя
func NotifyRegister(username string, email string) error {
	payload := map[string]interface{}{
		"notify_type": "user_register",
		"email":       email,
		"message":     username + " зарегистрирован",
	}

	return sendPayload(payload)
}
