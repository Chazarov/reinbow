package handlers

// Пакет handlers используется для реализации функций, которые должны вызываться
// при достижении endpoint'ов.

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
	"traineesheep/notifyservice/email"
	"traineesheep/notifyservice/errs"
	"traineesheep/notifyservice/tgbot"
	"traineesheep/notifyservice/webhook_handler"

	"github.com/go-telegram/bot"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Структура DTO используется для передачи данных из главного main.go файла
type DTO struct {
	bot            *bot.Bot        // Сущность телеграм бота
	sql_connection *pgxpool.Pool   // Коннект к базе
	smtp           *email.SmtpDTO  // Структура SMTP
	grtsChannels   chan int        // Канал нужный для контроля количества запросов
	wg             *sync.WaitGroup // WaitGroup для отметки, что задача выполнилась
	JwtSecret      string          // Секрет для JWT
}

// NewDTO возвращает новый экземпляр структуры DTO
// Используется для создания экземпляра структуры DTO
func NewDTO(bot *bot.Bot, conn *pgxpool.Pool, smtp *email.SmtpDTO, channel chan int, wg *sync.WaitGroup, secret string) *DTO {
	return &DTO{bot: bot, sql_connection: conn, smtp: smtp, grtsChannels: channel, wg: wg, JwtSecret: secret}
}

// Переменная для создания WebSocket соединения
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Массив с допустимыми типами уведомлений
var notify_types_allowed = []string{"user_register", "user_login", "admin_newImg",
	"user_imgVerdict", "user_email_confirmation"}

// Структура ResponseData используется для формирования ответа на запрос
type ResponseData struct {
	Success       bool   `json:"success"`     // Удачный ли запрос
	Error_message string `json:"err_message"` // Сообщение об ошибке
	Token         string `json:"token"`       // Токен
}

// Структура Recipent содержит данные, которые будут использоваться для отправки
// уведомления пользователю
type Recipent struct {
	Email       string `json:"email"`       // Почта пользователя
	Notify_Type string `json:"notify_type"` // Тип уведомления
	Message     string `json:"message"`     // Сообщение
	Telegram    int64  `json:"telegram_id"` // Телеграм ID пользователя
}

// Структура NotifyTypeMessenger сдержит данные, куда отправлять конкретный тип уведомлений
type NotifyTypeMessenger struct {
	NotifyType   string   `json:"notify_type"`        // Тип уведмоления
	Description  string   `json:"notify_description"` // Описание уведомления
	WantEmail    bool     `json:"want_email"`         // Отправляем ли на почту
	WantTelegram bool     `json:"want_telegram"`      // Отправляем ли в телеграм
	WebhookUrls  []string `json:"webhook_urls"`       // массив URL
}

// Структура NotifyTypeMessengerList нужна для хранения данных, полученных с сайта
// настройки уведомлений
type NotifyTypeMessengerList struct {
	// Success     bool                  `json:"Success"`     // Успех или нет
	Data []NotifyTypeMessenger `json:"data"` // Массив, содержащий данные куда отправлять определенные уведомления
}

// Структура LoginData нужна для хранения данных, передаваемых с сайта при авторизации
type LoginData struct {
	Login    string `json:"login"`    // Логин пользователя
	Password string `json:"password"` // Пароль пользователя
}

// Функция EnableCors настраивает CORS, чтобы можно было принимать запросы из браузера
func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

}

// Функция HandleNotify используется для отправки уведомлений
// На вход приходит запрос, с указанным типом уведомления и сообщением
// На выходе отправляется статус код и ответ, успешна ли была отправка уведомлений(-я)
func (d DTO) HandleNotify(w http.ResponseWriter, r *http.Request) {
	var response ResponseData

	enableCors(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Ограничиваем количество одновременно запущенных горутин
	// Чтобы обрабатывать одновременно за раз не больше CHANNEL_SIZE
	d.wg.Add(1)
	defer d.wg.Done()

	grtChan := d.grtsChannels
	grtChan <- 1
	defer func() { <-grtChan }()

	sql_conn := d.sql_connection

	var req Recipent

	body_data, err := io.ReadAll(r.Body)
	if err != nil {
		response.Success = false
		response.Error_message = errs.ErrReadingRequestMessage + err.Error()
		response_byte, _ := json.MarshalIndent(response, "", "    ")
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write(response_byte); err != nil {
			fmt.Println(errs.ErrWritingToRespBody)
		}
		return
	}

	if err := json.Unmarshal(body_data, &req); err != nil {
		response.Success = false
		response.Error_message = errs.ErrJsonUnmarshal + err.Error()
		response_byte, _ := json.MarshalIndent(response, "", "    ")
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write(response_byte); err != nil {
			fmt.Println(errs.ErrWritingToRespBody)
		}
		return
	}

	var tg_user_id int64 = req.Telegram
	var user_email string = req.Email

	var isAllowed bool = false

	for _, elem := range notify_types_allowed {
		if elem == req.Notify_Type {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		response.Success = false
		response.Error_message = "Недопустимый тип уведомления!"
		response_byte, _ := json.MarshalIndent(response, "", "    ")
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write(response_byte); err != nil {
			fmt.Println(errs.ErrWritingToRespBody)
			return
		}
	}

	// Если пользователь зарегался, то добавляю в свою базу его почту и user_id
	if req.Notify_Type == "user_register" {
		query := `INSERT INTO client (email, telegram_id) VALUES ($1, $2)`
		if _, err := sql_conn.Exec(context.Background(), query, user_email, tg_user_id); err != nil {
			fmt.Println("err:", err)
			response.Success = false
			response.Error_message = "При попытке добавить данные в базу произошла ошибка: " + err.Error() + err.Error()
			response_byte, _ := json.MarshalIndent(response, "", "    ")
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write(response_byte); err != nil {
				fmt.Println(errs.ErrWritingToRespBody)
				return
			}
			return
		}
	}
	if req.Notify_Type == "user_email_confirmation" {
		var emails_to_send []string
		notification_message := req.Message
		emails_to_send = append(emails_to_send, user_email)
		d.smtp.SendMessage(emails_to_send, []byte(notification_message), req.Notify_Type)
		response.Success = true
		response.Error_message = ""
		response_byte, _ := json.MarshalIndent(response, "", "    ")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(response_byte); err != nil {
			fmt.Println(errs.ErrWritingToRespBody)
			return
		}
		return
	}

	var want_email bool
	var want_telegram bool
	var want_webhook []string

	emails_to_send := make([]string, 0)

	// Получаем значения переключателей куда мы хотим получить уведомление
	query := `SELECT want_email, want_telegram, want_webhook FROM notify_type_message
	WHERE notify_type = $1`

	row := sql_conn.QueryRow(context.Background(), query, req.Notify_Type)
	if err := row.Scan(&want_email, &want_telegram, &want_webhook); err != nil {
		fmt.Println("Error while scanning db:", err)
		return
	}

	notification_message := req.Message

	fmt.Println(want_telegram, want_email, want_webhook)

	// Если мы хотим уведомление в ТГ
	if want_telegram {
		tgbot.HandleSendMessage(d.bot, context.Background(), tg_user_id, notification_message)
		fmt.Println("Я отправил сообщение в ТГ")
	}

	// Если мы хотим уведомление по Email
	if want_email {
		emails_to_send = append(emails_to_send, user_email)
		d.smtp.SendMessage(emails_to_send, []byte(notification_message), req.Notify_Type)
		fmt.Println("Отправил сообещние по почте")
	}

	// Если мы хотим уведомление по Webhook
	if len(want_webhook) != 0 {
		for _, url := range want_webhook {
			webhook_handler.SendWebhookMessage(url, []byte(notification_message))
		}
		fmt.Println("Отправил сообщение по вебхуку")
	}

	response.Success = true
	response.Error_message = ""
	response_byte, _ := json.MarshalIndent(response, "", "    ")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response_byte); err != nil {
		fmt.Println(errs.ErrWritingToRespBody)
		return
	}
	log.Printf("Уведомление типа %v было успешно доставлено\n", req.Notify_Type)
}

// Функция для сохранения значений переключателей, куда мы хотим получать уведомления
// На вход получаем запрос с значениями переключателей
// На выходе отправляется статус код и ответ, успешно ли было сохранение
func (d DTO) HandleSaveSettingsCheckmarks(w http.ResponseWriter, r *http.Request) {
	var response ResponseData

	enableCors(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	_, err := ValidateToken(r, d.JwtSecret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response.Success = false
		response.Error_message = "Unauthorized"
		respBytes, _ := json.MarshalIndent(response, "", "    ")
		if _, err := w.Write(respBytes); err != nil {
			fmt.Println("Failed to write response:", err)
			return
		}
		return
	}

	sql_conn := d.sql_connection

	body_byte, err := io.ReadAll(r.Body)
	if err != nil {
		response.Success = false
		response.Error_message = errs.ErrReadingRequestMessage + err.Error()
		respBytes, _ := json.MarshalIndent(response, "", "    ")
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write(respBytes); err != nil {
			fmt.Println(errs.ErrWritingToRespBody)
			return
		}
		return
	}

	fmt.Println("Я получил...\n", string(body_byte), "\n")

	var json_list NotifyTypeMessengerList
	if err := json.Unmarshal(body_byte, &json_list); err != nil {
		response.Success = false
		response.Error_message = errs.ErrJsonUnmarshal + err.Error()
		respBytes, _ := json.MarshalIndent(response, "", "    ")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		if _, err := w.Write(respBytes); err != nil {
			fmt.Println(errs.ErrWritingToRespBody)
			return
		}
		return
	}

	fmt.Println("Я успешно закинул в json_list")

	// После того как получили данные записываем их в базу, не создавая дубликаты
	for _, elem := range json_list.Data {
		fmt.Println("Читаю:", elem)
		query := `INSERT INTO notify_type_message (notify_type, notify_description, want_telegram, want_email, want_webhook)
    		VALUES ($1, $2, $3, $4, $5)
    		ON CONFLICT (notify_type) DO UPDATE SET
			notify_description = EXCLUDED.notify_description,
        	want_telegram = EXCLUDED.want_telegram,
        	want_email = EXCLUDED.want_email,
			want_webhook = EXCLUDED.want_webhook;`

		_, err := sql_conn.Exec(context.Background(), query, elem.NotifyType, elem.Description, elem.WantTelegram, elem.WantEmail, elem.WebhookUrls)
		if err != nil {
			response.Success = false
			response.Error_message = "Ошибка при вставке данных в базу. " + err.Error()
			respBytes, _ := json.MarshalIndent(response, "", "    ")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println(response.Error_message)
			if _, err := w.Write(respBytes); err != nil {
				fmt.Println(errs.ErrWritingToRespBody)
			}
			return
		}
	}

	response.Success = true
	response.Error_message = ""
	respBytes, _ := json.MarshalIndent(response, "", "    ")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(respBytes); err != nil {
		fmt.Println(errs.ErrWritingToRespBody)
		return
	}
	fmt.Println("Я успешно записал настройки")
}

// Функция для получения настроек перключателей куда мы хотим получать уведомление
// (нужно для того чтобы при обновлении страницы модератора восстанавливались значения
// переключателей которые уже сейчас лежат в базе)
func (d DTO) HandleGetNotifySettings(w http.ResponseWriter, r *http.Request) {
	var response ResponseData
	enableCors(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	_, err := ValidateToken(r, d.JwtSecret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response.Success = false
		response.Error_message = "Unauthorized"
		respBytes, _ := json.MarshalIndent(response, "", "    ")
		if _, err := w.Write(respBytes); err != nil {
			fmt.Println("Failed to write response:", err)
			return
		}
		return
	}

	var json_list = make([]NotifyTypeMessenger, 0)
	sql_conn := d.sql_connection

	query := `SELECT notify_type, notify_description, want_telegram, want_email, want_webhook FROM notify_type_message`
	rows, err := sql_conn.Query(context.Background(), query)
	if err != nil {
		response.Success = false
		response.Error_message = "Ошибка при выборе данных из базы. " + err.Error()
		respBytes, _ := json.MarshalIndent(response, "", "    ")
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write(respBytes); err != nil {
			fmt.Println(errs.ErrWritingToRespBody)
			return
		}
		return
	}
	defer rows.Close()

	for rows.Next() {
		var notify_type string
		var notify_description string
		var want_email bool
		var want_telegram bool
		var webhook_urls []string
		if err := rows.Scan(&notify_type, &notify_description, &want_telegram, &want_email, &webhook_urls); err != nil {
			fmt.Println(err)
			response.Success = false
			response.Error_message = "Ошибка при обходе базы данных. " + err.Error()
			respBytes, _ := json.MarshalIndent(response, "", "    ")
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write(respBytes); err != nil {
				fmt.Println(errs.ErrWritingToRespBody)
			}
			return
		}
		json_list = append(json_list, NotifyTypeMessenger{
			NotifyType:   notify_type,
			Description:  notify_description,
			WantEmail:    want_email,
			WantTelegram: want_telegram,
			WebhookUrls:  webhook_urls,
		})
	}

	var json_data_list = NotifyTypeMessengerList{
		Data: json_list,
	}

	resp_byte, err := json.MarshalIndent(json_data_list, "", "    ")
	if err != nil {
		response.Success = false
		response.Error_message = errs.ErrJsonMarshal + err.Error()
		respBytes, _ := json.MarshalIndent(response, "", "    ")
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write(respBytes); err != nil {
			fmt.Println(errs.ErrWritingToRespBody)
			return
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp_byte); err != nil {
		fmt.Println(errs.ErrWritingToRespBody)
		return
	}
	fmt.Println("Я успешно получил настройки")
}

// Функция HandleModeratorLogin используется для проверки данных для авторизации модератора
// На вход она получает в JSON логин и пароль потенциального пользователя
// На выход дает ответ, правильные ли данные были введены
func (d DTO) HandleModeratorLogin(w http.ResponseWriter, r *http.Request) {
	var response ResponseData
	var login_data LoginData
	var (
		jwt_key []byte
		jwt_t   *jwt.Token
	)

	enableCors(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	body_byte, err := io.ReadAll(r.Body)
	if err != nil {
		response.Success = false
		response.Error_message = errs.ErrReadingRequestMessage + err.Error()
		response_byte, _ := json.MarshalIndent(response, "", "    ")
		fmt.Println(response.Error_message)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write(response_byte); err != nil {
			fmt.Println(errs.ErrWritingToRespBody)
		}
		return
	}

	if err := json.Unmarshal(body_byte, &login_data); err != nil {

		var response ResponseData
		response.Success = false
		response.Error_message = errs.ErrJsonUnmarshal + ": " + err.Error()
		response_byte, _ := json.MarshalIndent(response, "", "    ")
		fmt.Println(response.Error_message)
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write(response_byte); err != nil {
			fmt.Println(errs.ErrWritingToRespBody)
		}
		return
	}

	//fmt.Println(login_data.Login, login_data.Password)

	if login_data.Login == "admin" && login_data.Password == "12345" {

		jwt_key = []byte(d.JwtSecret)
		jwt_t = jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{
				"iss": "my_auth_server",
				"sub": "admin",
				"exp": time.Now().Add(15 * time.Minute).Unix(),
			})
		jwt_s, err := jwt_t.SignedString(jwt_key)
		if err != nil {
			fmt.Println("Ошибка при создании криптографической подписи для токена", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var response ResponseData
		response.Success = true
		response.Error_message = ""
		response.Token = jwt_s
		response_byte, _ := json.MarshalIndent(response, "", "    ")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(response_byte); err != nil {
			fmt.Println(errs.ErrWritingToRespBody)
		}
		fmt.Println("Логин и пароль успешно прошли проверку!")
		return
	} else {
		var response ResponseData
		response.Success = false
		response.Error_message = "Неправильно введен логин или пароль. Попробуйте еще раз!"
		response_byte, _ := json.MarshalIndent(response, "", "    ")
		fmt.Println(response.Error_message)
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(response_byte); err != nil {
			fmt.Println(errs.ErrWritingToRespBody)
		}
		fmt.Println("Логин и пароль не прошли проверку...")
		return
	}
}

func ValidateToken(r *http.Request, jwtSecret string) (jwt.MapClaims, error) {

	fmt.Println("Я получил токен:", jwtSecret)

	authHeader := r.Header.Get("Authorization")

	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT секрет не задан")
	}

	if authHeader == "" {
		err := "Проверка токена не прошла! Необходима авторизация!"
		fmt.Println(err)
		return nil, fmt.Errorf(err)
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		err := "Неправильный формат хедера: Необходимо использовать Bearer <token>"
		fmt.Println(err)
		return nil, fmt.Errorf(err)
	}

	tokenString := headerParts[1]

	fmt.Println("Получил token string:", tokenString)

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})

	fmt.Println(err, "Валиден ли?", token.Valid)

	if err != nil || !token.Valid {
		err := "Неправильный токен или токен истёк!"
		fmt.Println(err)
		return nil, fmt.Errorf(err)
	}

	return claims, nil
}
