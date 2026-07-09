package service

import (
	"errors"
	"strings"

	"traineesheep/feedservice/internal/client/notify_service"
	"traineesheep/feedservice/internal/repository"

	model "github.com/tokyobordel/traineepkg/models"
)

// UserService предоставляет методы для работы с пользователями:
// создание, проверка существования, получение по имени или ID.
type UserService struct {
	UserDAO      *repository.UserDAO
	NotifyClient *client.NotifyClient
}

type registerInput struct {
	Email string
}

// parseRegisterData извлекает и валидирует данные из map
func parseRegisterData(data map[string]string) (registerInput, error) {
	input := registerInput{}

	// Email обязателен
	email, ok := data["email"]
	if !ok || strings.TrimSpace(email) == "" {
		return input, errors.New("email is required")
	}
	input.Email = email

	return input, nil
}

// NewUserService создаёт новый экземпляр UserService с переданным DAO.
func NewUserService(userDAO *repository.UserDAO, notifyClient *client.NotifyClient) *UserService {
	return &UserService{UserDAO: userDAO, NotifyClient: notifyClient}
}

// ExistsByUsername проверяет, существует ли пользователь с указанным именем.
// Возвращает true, если пользователь найден, и ошибку при проблемах с БД.
func (us *UserService) ExistsByUsername(username string) (bool, error) {
	return us.UserDAO.ExistsByUsername(username)
}

// ExistsByUsername проверяет, существует ли пользователь с указанным email.
// Возвращает true, если пользователь найден, и ошибку при проблемах с БД.
func (us *UserService) ExistsByEmail(email string) (bool, error) {
	return us.UserDAO.ExistsByEmail(email)
}

// GetByUsername возвращает модель пользователя по его имени.
// Если пользователь не найден, возвращается пустая структура и ошибка sql.ErrNoRows.
func (us *UserService) GetByUsername(username string) (model.User, error) {
	return us.UserDAO.GetByUsername(username)
}

// ConfirmUserAccount подтверждает регистрацию пользователя.
func (us *UserService) ConfirmUserAccount(userID int) error {
	return us.UserDAO.ConfirmUserAccount(userID)
}
