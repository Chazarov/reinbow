package repository

import (
	"database/sql"
	models "traineesheep/feedservice/internal/model"
)

// TokenDAO обеспечивает доступ к операциям с refresh-токенами в базе данных.
// В текущей итерации хранение токенов реализовано в памяти, поэтому
// некоторые методы являются заглушками.
type TokenDAO struct {
	db *sql.DB
}

// NewTokenDAO создаёт новый экземпляр TokenDAO с заданным подключением к БД.
func NewTokenDAO(db *sql.DB) *TokenDAO {
	return &TokenDAO{db: db}
}

// DeleteToken удаляет переданный refresh-токен из базы данных.
// В текущей версии приложения инвалидация токенов не реализована
// поэтому метод всегда возвращает nil.
func (td *TokenDAO) DeleteToken(token string) error {
	return nil
}

// GetRefreshToken извлекает из БД запись refresh-токена по его строковому
// представлению. Возвращает модель RefreshToken и ошибку, если токен не найден
// или запрос выполнить не удалось.
func (td *TokenDAO) GetRefreshToken(token string) (models.RefreshToken, error) {
	var rt models.RefreshToken
	err := td.db.QueryRow(
		"SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token = $1",
		token,
	).Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt)
	return rt, err
}
