package repository

import (
	"database/sql"
	"time"
	models "traineesheep/feedservice/internal/model"
	"traineesheep/feedservice/internal/utils"

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
	// Хэшируем пароль
	hashedPassword, passError := bcrypt.GenerateFromPassword([]byte(userData.Password),
		bcrypt.DefaultCost)
	if passError != nil {
		return models.User{}, passError
	}

	var user models.User
	dbError := ud.db.QueryRow(
		"INSERT INTO users (username, password, created_at, email) "+
			"VALUES ($1, $2, $3, $4) RETURNING id, username, created_at, email",
		userData.Username, string(hashedPassword), time.Now(), userData.Email,
	).Scan(&user.ID, &user.Username, &user.CreatedAt, &user.Email)
	if dbError != nil {
		return models.User{}, dbError
	}
	return user, nil
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

// GetByUsername возвращает модель пользователя по его имени.
// Если пользователь не найден, возвращается пустая структура и ошибка sql.ErrNoRows.
func (ud *UserDAO) GetByUsername(username string) (models.User, error) {
	var user models.User
	err := ud.db.QueryRow(
		"SELECT id, username, password, created_at FROM users WHERE username = $1",
		username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)

	return user, err
}

// GetByID возвращает модель пользователя по его идентификатору.
// Если пользователь не найден, возвращается пустая структура и ошибка sql.ErrNoRows.
func (ud *UserDAO) GetByID(userID int) (models.User, error) {
	var user models.User
	err := ud.db.QueryRow(
		"SELECT id, username, password, created_at FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)

	return user, err
}
