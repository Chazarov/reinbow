package service

import (
	"traineesheep/feedservice/internal/model"
	"traineesheep/feedservice/internal/repository"
	"traineesheep/feedservice/internal/utils"
)

// UserService предоставляет методы для работы с пользователями:
// создание, проверка существования, получение по имени или ID.
type UserService struct {
	userDAO *repository.UserDAO
}

// NewUserService создаёт новый экземпляр UserService с переданным DAO.
func NewUserService(userDAO *repository.UserDAO) *UserService {
	return &UserService{userDAO: userDAO}
}

// CreateUser создаёт нового пользователя на основе переданных данных.
// Пароль внутри DAO хэшируется с помощью bcrypt.
// Возвращает созданную модель пользователя и ошибку.
func (us *UserService) CreateUser(user utils.UserData) (models.User, error) {
	return us.userDAO.CreateUser(user)
}

// ExistsByUsername проверяет, существует ли пользователь с указанным именем.
// Возвращает true, если пользователь найден, и ошибку при проблемах с БД.
func (us *UserService) ExistsByUsername(username string) (bool, error) {
	return us.userDAO.ExistsByUsername(username)
}

// GetByUsername возвращает модель пользователя по его имени.
// Если пользователь не найден, возвращается пустая структура и ошибка sql.ErrNoRows.
func (us *UserService) GetByUsername(username string) (models.User, error) {
	return us.userDAO.GetByUsername(username)
}

// GetByID возвращает модель пользователя по его идентификатору.
// Если пользователь не найден, возвращается пустая структура и ошибка sql.ErrNoRows.
func (us *UserService) GetByID(userID int) (models.User, error) {
	return us.userDAO.GetByID(userID)
}
