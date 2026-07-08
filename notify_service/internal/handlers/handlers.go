// Package handlers предоставляет обработчики HTTP-запросов для сервиса уведомлений.
//
// Данный пакет содержит реализацию всех endpoint'ов API сервиса,
// включая обработку уведомлений, настройку типов уведомлений,
// получение настроек и авторизацию модераторов.
package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"traineesheep/notifyservice/pkg/email"

	"github.com/go-telegram/bot"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DTO (Data Transfer Object) представляет собой структуру данных,
// которая инкапсулирует все зависимости и бизнес-логику для обработчиков HTTP-запросов.
//
// Структура содержит:
//   - Telegram бот для отправки сообщений
//   - Пул соединений с базой данных PostgreSQL
//   - Конфигурацию SMTP для отправки email уведомлений
//   - Канал для передачи данных между goroutines
//   - WaitGroup для синхронизации горутин
//   - Секретный ключ для JWT токенов
//   - Логин от аккаунта в админ панели
//   - Пароль от аккаунта в админ панели
type DTO struct {
	bot            *bot.Bot
	sql_connection *pgxpool.Pool
	smtp           *email.SmtpDTO
	grtsChannels   chan int
	wg             *sync.WaitGroup
	JwtSecret      string
	admin_login    string
	admin_pass     string
}

// NewDTO создает новый экземпляр DTO с заданными параметрами.
//
// Параметры:
//   - bot: экземпляр Telegram бота
//   - dbpool: пул соединений с PostgreSQL базой данных
//   - smtpDto: конфигурация SMTP сервера
//   - grtChan: канал для передачи данных между goroutines
//   - wg: wait group для синхронизации горутин
//   - jwtSecret: секретный ключ для JWT токенов
//
// Возвращает указатель на созданный DTO.
func NewDTO(bot *bot.Bot, conn *pgxpool.Pool, smtp *email.SmtpDTO, channel chan int, wg *sync.WaitGroup, secret string, login string, pass string) *DTO {
	return &DTO{
		bot:            bot,
		sql_connection: conn,
		smtp:           smtp,
		grtsChannels:   channel,
		wg:             wg,
		JwtSecret:      secret,
		admin_login:    login,
		admin_pass:     pass,
	}
}

// Переменная для создания WebSocket соединения
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Массив с допустимыми типами уведомлений
var notify_types_allowed = []string{}

// Функция EnableCors настраивает CORS, чтобы можно было принимать запросы из браузера
func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

}

// Функция ValidateToken нужна для валидация JWT токена
func ValidateToken(r *http.Request, jwtSecret string) (jwt.MapClaims, error) {
	logMessage += "Был получен запрос на проверку JWT токена\n"
	logMessage += "=========================\n"

	authHeader := r.Header.Get("Authorization")

	if jwtSecret == "" {
		logMessage += "JWT секрет оказался не задан...\n"
		return nil, fmt.Errorf("JWT секрет не задан")
	}

	if authHeader == "" {
		err := "Проверка токена не прошла! Необходима авторизация!"
		logMessage += err + "\n"
		return nil, fmt.Errorf(err)
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		err := "Неправильный формат хедера: Необходимо использовать Bearer <token>"
		logMessage += err + "\n"
		return nil, fmt.Errorf(err)
	}

	tokenString := headerParts[1]

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		err := "Неправильный токен или токен истёк!"
		logMessage += err + "\n"
		return nil, fmt.Errorf(err)
	}

	logMessage += "Проверка токена прошла успешно!\n"
	log.Println(logMessage)
	return claims, nil
}
