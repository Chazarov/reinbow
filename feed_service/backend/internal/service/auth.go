package service

import (
	"traineesheep/feedservice/internal/middleware"

	authService "github.com/tokyobordel/traineepkg/auth/service"
	model "github.com/tokyobordel/traineepkg/models"
	"golang.org/x/crypto/bcrypt"

	"log"
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

func (us *AuthService) Login(pass string, login string) (model.User, error) {
	hash, err := us.UserDAO.GetPasswordHash(login)
	if err != nil {
		return model.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		return model.User{}, err
	}
	return us.UserDAO.GetByUsername(login)
}

// Register создаёт нового пользователя на основе переданных данных
// и отправляет уведомление о регистрации
// Возвращает созданную модель пользователя и ошибку.
func (us *AuthService) Register(pass string, login string, data map[string]string) (model.User, error) {
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
	user, err := us.UserDAO.CreateUser(userData)
	if err != nil {
		return model.User{}, err
	}
	err = us.NotifyClient.NotifyRegisterForAdmin(userData.Username, userData.Email)
	if err != nil {
		log.Println(err.Error())
	}

	token, err := middleware.GenerateConfirmToken(user.ID)
	if err != nil {
		return model.User{}, err
	}
	us.NotifyClient.SendUserdataEmail(userData.Username, userData.Email, token)

	return user, err
}

// GetMe возвращает модель пользователя по его идентификатору.
// Если пользователь не найден, возвращается пустая структура и ошибка sql.ErrNoRows.
func (us *AuthService) GetMe(id int) (model.User, error) {
	return us.UserDAO.GetByID(id)
}
