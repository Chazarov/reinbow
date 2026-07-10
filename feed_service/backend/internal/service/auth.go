package service

import (
	"fmt"

	authService "github.com/tokyobordel/traineepkg/auth/service"
	model "github.com/tokyobordel/traineepkg/models"
	"golang.org/x/crypto/bcrypt"

	"traineesheep/feedservice/internal/client/notify_service"
	"traineesheep/feedservice/internal/repository"
	"traineesheep/feedservice/internal/utils"
)

// UserService предоставляет методы для работы с пользователями:
// создание, проверка существования, получение по имени или ID.
type AuthService struct {
	UserDAO      *repository.UserDAO
	NotifyClient *client.NotifyClient
}

// NewUserService создаёт новый экземпляр UserService с переданным DAO.
func NewAuthService(userDAO *repository.UserDAO, notifyClient *client.NotifyClient) authService.IAuthService {
	return &AuthService{UserDAO: userDAO, NotifyClient: notifyClient}
}

func (authService *AuthService) Login(pass string, login string) (model.User, error) {
	hash, err := authService.UserDAO.GetPasswordHash(login)
	if err != nil {
		return model.User{}, fmt.Errorf("Неверный логин или пароль")
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		return model.User{}, fmt.Errorf("Неверный логин или пароль")
	}
	return authService.UserDAO.GetByUsername(login)
}

// Register создаёт нового пользователя на основе переданных данных
// и отправляет уведомление о регистрации
// Возвращает созданную модель пользователя и ошибку.
func (authService *AuthService) Register(pass string, login string, data map[string]string) (model.User, error) {
	input, err := parseRegisterData(data)
	if err != nil {
		return model.User{}, err
	}

	// Формируем UserData для DAO
	userData := utils.UserData{
		Username: login, // логин = username
		Password: pass,
		Email:    input.Email,
	}
	user, err := authService.UserDAO.CreateUser(userData)
	if err != nil {
		return model.User{}, fmt.Errorf("Пользователь с указанными данными уже существует")
	}

	token, err := utils.GenerateConfirmToken(user.ID)
	if err != nil {
		return model.User{}, fmt.Errorf("Ошибка создания токена для потдверждения почты. Попробуйте позже")
	}
	authService.NotifyClient.SendUserdataEmail(userData.Username, userData.Email, token)

	return user, err
}

// GetMe возвращает модель пользователя по его идентификатору.
// Если пользователь не найден, возвращается пустая структура и ошибка sql.ErrNoRows.
func (authService *AuthService) GetMe(id int) (model.User, error) {
	return authService.UserDAO.GetByID(id)
}
