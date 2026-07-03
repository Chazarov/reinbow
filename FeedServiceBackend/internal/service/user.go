package service

import (
	"traineesheep/feedservice/internal/model"
	"traineesheep/feedservice/internal/repository"
	"traineesheep/feedservice/internal/utils"
)

type UserService struct {
	userDAO *repository.UserDAO
}

func NewUserService(userDAO *repository.UserDAO) *UserService {
	return &UserService{userDAO: userDAO}
}

func (us *UserService) CreateUser(user utils.UserData) (models.User, error) {
	return us.userDAO.CreateUser(user)
}

func (us *UserService) ExistsByUsername(username string) (bool, error) {
	return us.userDAO.ExistsByUsername(username)
}

func (us *UserService) GetByUsername(username string) (models.User, error) {
	return us.userDAO.GetByUsername(username)
}

func (us *UserService) GetByID(userID int) (models.User, error) {
	return us.userDAO.GetByID(userID)
}
