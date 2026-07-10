// Пакет service содержит доменную бизнес-логику сервиса изображений.
package service

import (
	"context"
	"strconv"
	"time"

	"traineesheep/imageservice/internal/app/models"
	"traineesheep/imageservice/internal/app/ports/repository"

	authContract "github.com/tokyobordel/traineepkg/auth/service"
	"github.com/tokyobordel/traineepkg/errors"
	"github.com/tokyobordel/traineepkg/logger"
	pkgModels "github.com/tokyobordel/traineepkg/models"

	"golang.org/x/crypto/bcrypt"
)

// AuthService реализует аутентификацию клиентов через traineepkg.
type AuthService struct {
	repo   *repository.ClientPostgresRepository
	logger *logger.ContextLogger
}

// NewPkgAuthService создаёт реализацию сервиса аутентификации для traineepkg.
func NewPkgAuthService(repo *repository.ClientPostgresRepository, logger *logger.ContextLogger) authContract.IAuthService {
	return &AuthService{
		repo:   repo,
		logger: logger,
	}
}

// hashPassword возвращает bcrypt-хеш переданного пароля.
func (s *AuthService) hashPassword(password string) (string, errors.DomainError) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.NewInternalServiceError("Error when trying to create a password hash", err)
	}
	return string(hashed), nil
}

// Login проверяет учётные данные и возвращает профиль пользователя.
func (s *AuthService) Login(pass string, login string) (pkgModels.User, error) {
	ctx := context.Background()
	if login == "" {
		return pkgModels.User{}, errors.NewInvalidParametersError("login", login, "login cannot be empty")
	}
	if pass == "" {
		return pkgModels.User{}, errors.NewInvalidParametersError("pass", pass, "pass cannot be empty")
	}

	client, err := s.repo.GetClientByName(ctx, login)
	if err != nil {
		return pkgModels.User{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(client.PassHash), []byte(pass)); err != nil {
		s.logger.Errorf(ctx, "Invalid credentials for login %s", login)
		return pkgModels.User{}, errors.NewUnauthorizedError("invalid credentials")
	}

	return clientToPkgUser(client), nil
}

// Register создаёт нового пользователя и возвращает его профиль.
func (s *AuthService) Register(pass string, login string, data map[string]string) (pkgModels.User, error) {
	ctx := context.Background()
	if login == "" {
		return pkgModels.User{}, errors.NewInvalidParametersError("login", login, "login cannot be empty")
	}
	if pass == "" {
		return pkgModels.User{}, errors.NewInvalidParametersError("pass", pass, "pass cannot be empty")
	}

	existingClient, err := s.repo.GetClientByName(ctx, login)
	if err == nil && existingClient.Id > 0 {
		return pkgModels.User{}, errors.NewUniqueConstraintError("client", "name", login)
	}
	if err != nil && !errors.IsNotFoundError(err) {
		return pkgModels.User{}, errors.NewInternalServiceError("failed to validate user uniqueness", err)
	}

	passHash, err := s.hashPassword(pass)
	if err != nil {
		return pkgModels.User{}, err
	}

	client, err := s.repo.CreateClient(ctx, models.CreateClientMto{
		PassHash:  passHash,
		Name:      login,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return pkgModels.User{}, err
	}

	user := clientToPkgUser(client)
	if data != nil {
		if user.Data == nil {
			user.Data = make(map[string]string)
		}
		for key, value := range data {
			user.Data[key] = value
		}
	}

	return user, nil
}

// GetMe возвращает профиль пользователя по идентификатору.
func (s *AuthService) GetMe(id int) (pkgModels.User, error) {
	ctx := context.Background()
	if id <= 0 {
		return pkgModels.User{}, errors.NewInvalidParametersError("id", id, "id must be positive")
	}

	client, err := s.repo.GetClientByID(ctx, strconv.Itoa(id))
	if err != nil {
		return pkgModels.User{}, err
	}

	return clientToPkgUser(client), nil
}

// clientToPkgUser преобразует клиента доменной модели в модель traineepkg.
func clientToPkgUser(client models.Client) pkgModels.User {
	return pkgModels.User{
		ID:    client.Id,
		Login: client.Name,
		Data: map[string]string{
			"created_at": client.CreatedAt.Format(time.RFC3339),
		},
	}
}
