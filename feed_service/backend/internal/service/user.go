package service

import (
	"errors"
	"fmt"
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
func (userService *UserService) ExistsByUsername(username string) (bool, error) {
	return userService.UserDAO.ExistsByUsername(username)
}

// ExistsByUsername проверяет, существует ли пользователь с указанным email.
// Возвращает true, если пользователь найден, и ошибку при проблемах с БД.
func (userService *UserService) ExistsByEmail(email string) (bool, error) {
	return userService.UserDAO.ExistsByEmail(email)
}

// GetByID возвращает модель пользователя по его ID.
// Если пользователь не найден, возвращается пустая структура и ошибка sql.ErrNoRows.
func (userService *UserService) GetByID(id int) (model.User, error) {
	return userService.UserDAO.GetByID(id)
}

// GetByUsername возвращает модель пользователя по его имени.
// Если пользователь не найден, возвращается пустая структура и ошибка sql.ErrNoRows.
func (userService *UserService) GetByUsername(username string) (model.User, error) {
	return userService.UserDAO.GetByUsername(username)
}

// ConfirmUserAccount подтверждает регистрацию пользователя.
func (userService *UserService) ConfirmUserAccount(userID int) error {
	return userService.UserDAO.ConfirmUserAccount(userID)
}

// ConfirmUserAccount подтверждает регистрацию пользователя.
func (userService *UserService) SendNotificationExternal(userID int) error {
	user, err := userService.GetByID(userID)
	if err != nil {
		return err
	}

	email, ok := user.Data["email"]
	if !ok {
		return fmt.Errorf("У пользователя отсутствует email")
	}

	return userService.NotifyClient.NotifyRegisterForAdmin(user.Login, email)
}
