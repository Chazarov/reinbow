package repository

import (
	"database/sql"
	"strconv"
	"time"
	"traineesheep/feedservice/internal/utils"

	"github.com/tokyobordel/traineepkg/models"

	"golang.org/x/crypto/bcrypt"
)

// UserDAO обеспечивает доступ к операциям с пользователями в базе данных.
type UserDAO struct {
	db *sql.DB
}

// NewUserDAO создаёт новый экземпляр UserDAO с заданным подключением к БД.
func NewUserDAO(db *sql.DB) *UserDAO {
	return &UserDAO{db: db}
}

// CreateUser создаёт нового пользователя с переданными данными.
// Пароль хэшируется с помощью bcrypt перед сохранением.
// Возвращает созданную модель User и ошибку.
func (ud *UserDAO) CreateUser(userData utils.UserData) (models.User, error) {
	var id int
	var username string
	var createdAt time.Time
	var email string

	hashedPassword, passError := bcrypt.GenerateFromPassword([]byte(userData.Password),
		bcrypt.DefaultCost)
	if passError != nil {
		return models.User{}, passError
	}

	dbError := ud.db.QueryRow(
		"INSERT INTO users (username, password, created_at, email) "+
			"VALUES ($1, $2, $3, $4) RETURNING id, username, created_at, email",
		userData.Username, string(hashedPassword), time.Now(), userData.Email,
	).Scan(&id, &username, &createdAt, &email)

	if dbError != nil {
		return models.User{}, dbError
	}
	// Собираем целевую модель
	result := models.User{
		ID:    id,
		Login: username, // сопоставляем логину имя пользователя
		Data: map[string]string{
			"email":        email,
			"is_confirmed": strconv.FormatBool(false),
			"created_at":   createdAt.Format(time.RFC3339),
		},
	}
	return result, nil
}

// ExistsByUsername проверяет, существует ли пользователь с указанным именем.
// Возвращает true, если пользователь найден, и ошибку запроса.
func (ud *UserDAO) ExistsByUsername(username string) (bool, error) {
	var exists bool
	err := ud.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)",
		username,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// ExistsByEmail проверяет, существует ли пользователь с указанным email.
// Возвращает true, если пользователь найден, и ошибку запроса.
func (ud *UserDAO) ExistsByEmail(email string) (bool, error) {
	var exists bool
	err := ud.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)",
		email,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// GetByUsername возвращает модель пользователя по его имени.
// Если пользователь не найден, возвращается пустая структура и ошибка sql.ErrNoRows.
func (ud *UserDAO) GetPasswordHash(username string) (string, error) {
	var hash string
	err := ud.db.QueryRow(
		"SELECT password FROM users WHERE username = $1",
		username,
	).Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}

// GetByUsername возвращает модель пользователя по его имени.
// Если пользователь не найден, возвращается пустая структура и ошибка sql.ErrNoRows.
func (ud *UserDAO) GetByUsername(username string) (models.User, error) {
	var id int
	var login string
	var password string
	var createdAt time.Time
	var email string
	var isConfirmed bool

	err := ud.db.QueryRow(
		"SELECT id, username, password, created_at, email, is_confirmed FROM users WHERE username = $1",
		username,
	).Scan(&id, &login, &password, &createdAt, &email, &isConfirmed)

	if err != nil {
		return models.User{}, err
	}
	// Собираем целевую модель
	result := models.User{
		ID:    id,
		Login: username, // сопоставляем логину имя пользователя
		Data: map[string]string{
			"email":        email,
			"is_confirmed": strconv.FormatBool(false),
			"created_at":   createdAt.Format(time.RFC3339),
		},
	}
	return result, nil
}

// GetByID возвращает модель пользователя по его идентификатору.
// Если пользователь не найден, возвращается пустая структура и ошибка sql.ErrNoRows.
func (ud *UserDAO) GetByID(userID int) (models.User, error) {
	var id int
	var login string
	var password string
	var createdAt time.Time
	var email string
	var isConfirmed bool

	err := ud.db.QueryRow(
		"SELECT id, username, password, created_at, is_confirmed, email FROM users WHERE id = $1",
		userID,
	).Scan(&id, &login, &password, &createdAt, &isConfirmed, &email)

	if err != nil {
		return models.User{}, err
	}
	// Собираем целевую модель
	result := models.User{
		ID:    id,
		Login: login,
		Data: map[string]string{
			"email":        email,
			"is_confirmed": strconv.FormatBool(isConfirmed),
			"created_at":   createdAt.Format(time.RFC3339),
		},
	}

	return result, nil
}

// ConfirmUserAccount подтверждает регистрацию пользователя.
func (ud *UserDAO) ConfirmUserAccount(userID int) error {
	err := ud.db.QueryRow(
		"UPDATE users SET is_confirmed = TRUE WHERE id = $1;",
		userID,
	).Err()

	return err
}
