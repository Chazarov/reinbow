// Package handlers предоставляет обработчики HTTP-запросов для сервиса уведомлений.
//
// Данный пакет содержит реализацию всех endpoint'ов API сервиса,
// включая обработку уведомлений, настройку типов уведомлений,
// получение настроек и авторизацию модераторов.
package handlers

import (
	"net/http"

	"github.com/go-telegram/bot"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

var (
	TgBot *bot.Bot
	// GrtChannels                      chan int
	// Wg                               *sync.WaitGroup
	AdminLogin, AdminPass, JwtSecret string
	TaskGroup                        *errgroup.Group
)

// Переменная для создания WebSocket соединения
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Массив с допустимыми типами уведомлений
var notifyTypesAllowed = []string{}

// Функция EnableCors настраивает CORS, чтобы можно было принимать запросы из браузера
func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
