package repository

import (
	"database/sql"
	models "traineesheep/feedservice/internal/model"
)

type TokenDAO struct {
	db *sql.DB
}

func NewTokenDAO(db *sql.DB) *TokenDAO {
	return &TokenDAO{db: db}
}

func (td *TokenDAO) DeleteToken(token string) error {
	// todo инвалидация токена. В текущей итерации проекта смысла в этом нет
	// _, err := td.db.Exec("DELETE FROM refresh_tokens WHERE token = $1", token)
	return nil
}

func (td *TokenDAO) GetRefreshToken(token string) (models.RefreshToken, error) {
	var rt models.RefreshToken
	err := td.db.QueryRow(
		"SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token = $1",
		token,
	).Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt)
	return rt, err
}
