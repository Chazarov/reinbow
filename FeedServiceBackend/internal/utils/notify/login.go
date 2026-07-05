package notify

// NotifyLogin уведомляет пользователя о том, что в его аккаунт совершен вход
func NotifyLogin(username string, email string) {
	payload := map[string]interface{}{
		"notify_type": "user_login",
		"email":       email,
		"message":     "В ваш аккаунт (" + username + ") совершен вход",
		"telegram_id": 924956695,
	}

	sendPayload(payload)
}
