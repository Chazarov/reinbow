package service

import (
	"traineesheep/feedservice/internal/model"
	"traineesheep/feedservice/internal/repository"
)

// TokenService предоставляет методы для работы с refresh-токенами:
// удаление (инвалидация) токена и получение информации о токене.
type TokenService struct {
	TokenDAO *repository.TokenDAO
}

// NewTokenService создаёт новый экземпляр TokenService с переданным DAO.
func NewTokenService(tokenDAO *repository.TokenDAO) *TokenService {
	return &TokenService{TokenDAO: tokenDAO}
}

// DeleteToken удаляет refresh-токен из базы данных (инвалидация).
// В текущей версии приложения операция не выполняется (метод-заглушка),
// всегда возвращает nil.
func (ts *TokenService) DeleteToken(token string) error {
	return ts.TokenDAO.DeleteToken(token)
}

// GetRefreshToken возвращает модель RefreshToken по его строковому представлению.
// Если токен не найден, возвращается ошибка sql.ErrNoRows.
func (ts *TokenService) GetRefreshToken(token string) (models.RefreshToken, error) {
	return ts.TokenDAO.GetRefreshToken(token)
}
