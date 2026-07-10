package services

import (
	"fmt"

	authServicePkg "github.com/tokyobordel/traineepkg/auth/service"
	"github.com/tokyobordel/traineepkg/models"
)

type UserService struct {
	login   string
	pass    string
	adminId int
}

func NewUserService(adminPass string, adminLogin string) authServicePkg.IAuthService {
	return &UserService{
		login:   adminLogin,
		pass:    adminPass,
		adminId: 1,
	}

}

func (s *UserService) Login(pass string, login string) (models.User, error) {
	if pass == s.pass && login == s.login {
		return models.User{ID: s.adminId, Login: login, Data: nil}, nil
	} else {
		return models.User{ID: 0}, fmt.Errorf(" Access denied")
	}
}
func (s *UserService) Register(pass string, login string, data map[string]string) (models.User, error) {
	return models.User{ID: 0}, fmt.Errorf(" Method not Implemented")
}
func (s *UserService) GetMe(id int) (models.User, error) {
	if id == s.adminId {
		return models.User{ID: s.adminId, Login: s.login, Data: nil}, nil
	} else {
		return models.User{ID: 0}, fmt.Errorf(" User Not found")
	}
}
