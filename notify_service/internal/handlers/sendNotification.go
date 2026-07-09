package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"traineesheep/notifyservice/internal/database"
	"traineesheep/notifyservice/internal/errs"
	"traineesheep/notifyservice/internal/tgbot"
	"traineesheep/notifyservice/internal/types"
	"traineesheep/notifyservice/internal/webhook_handler"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logMessage string

// HandleNotify обрабатывает POST-запросы на отправку уведомлений.
//
// Endpoint: POST /api/notify
//
// Тело запроса должно содержать JSON с данными уведомления.
// Функция проверяет JWT токен, извлекает данные из запроса
// и запускает горутину для отправки уведомлений через различные каналы.
//
// Возможные коды ответа:
//   - 200: уведомление успешно принято к отправке
//   - 400: неверный формат запроса
//   - 401: невалидный или отсутствующий JWT токен
//   - 500: внутренняя ошибка сервера

// HandleNotify godoc
// @Summary Отправка уведомления
// @Description Принимает запрос на отправку уведомления пользователю по указанному типу.
// @Description Поддерживает отправку по Email, Telegram и Webhook в зависимости от настроек.
// @Description Требует JWT-токен в заголовке Authorization (если не требуется, уберите @Security).
// @Tags Notifications
// @Accept json
// @Produce json
// @Param request body types.Recipent true "Данные для отправки уведомления"
// @Success 200 {object} types.ResponseData "Уведомление успешно отправлено"
// @Failure 400 {object} types.ResponseData "Неверный формат запроса или недопустимый тип уведомления"
// @Failure 401 {object} types.ResponseData "Неавторизован (отсутствует или невалидный JWT токен)"
// @Failure 500 {object} types.ResponseData "Внутренняя ошибка сервера"
// @Router /api/notify [post]
// @Security ApiKeyAuth
func (d DTO) HandleNotify(w http.ResponseWriter, r *http.Request) {
	logMessage = ""
	var response types.ResponseData
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	var req types.Recipent
	database_conn_dto := database.NewDatabaseDTO(d.sql_connection)

	enableCors(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Ограничиваем количество одновременно запущенных горутин
	// Чтобы обрабатывать одновременно за раз не больше WORKER_COUNT
	d.wg.Add(1)
	defer d.wg.Done()

	grtChan := d.grtsChannels
	grtChan <- 1
	defer func() { <-grtChan }()

	body_data, err := io.ReadAll(r.Body)
	if err != nil {
		response.Success = false
		response.Error_message = errs.ErrReadingRequestMessage + err.Error()
		response_byte, _ := json.MarshalIndent(response, "", "    ")
		w.WriteHeader(http.StatusInternalServerError)
		log.Error().Msg(response.Error_message)
		if _, err := w.Write(response_byte); err != nil {
			log.Error().Msg(errs.ErrWritingToRespBody)
		}
		return
	}

	if err := json.Unmarshal(body_data, &req); err != nil {
		response.Success = false
		response.Error_message = errs.ErrJsonUnmarshal + err.Error()
		response_byte, _ := json.MarshalIndent(response, "", "    ")
		log.Error().Msg(response.Error_message)
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write(response_byte); err != nil {
			log.Error().Msg(errs.ErrWritingToRespBody)
		}
		return
	}

	logMessage += fmt.Sprintf("Был получен запрос на тип %v\n", req.Notify_Type)
	logMessage += fmt.Sprintf("Сообщение для отправки: \"%v\"\n", req.Message)
	logMessage += "=========================\n"

	var tg_user_id int64 = 924956695 // ВРЕМЕННО ЗАХАРДКОЖЕН ID ДО РЕАЛИЗАЦИИ ПОЛУЧЕНИЯ ИЗ ЧАТА
	var user_email string = req.Email

	var isAllowed bool = false

	rows, err := database_conn_dto.GetNotifyTypes()
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	if len(notify_types_allowed) == 0 {
		for rows.Next() {
			var temp string
			rows.Scan(&temp)
			notify_types_allowed = append(notify_types_allowed, temp)
		}
	}

	for _, elem := range notify_types_allowed {
		if elem == req.Notify_Type {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		response.Success = false
		response.Error_message = fmt.Sprintf("Недопустимый тип уведомления! %v", req.Notify_Type)
		response_byte, _ := json.MarshalIndent(response, "", "    ")
		log.Error().Msg(response.Error_message)
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write(response_byte); err != nil {
			log.Error().Msg(errs.ErrWritingToRespBody)
			return
		}
		return
	}

	// Если пользователь зарегался, то добавляю в свою базу его почту и user_id и выхожу
	if req.Notify_Type == "user_register" {
		if err := database_conn_dto.AddEmail(user_email); err != nil {
			response.Success = false
			response.Error_message = "При попытке добавить данные в базу произошла ошибка: " + err.Error()
			response_byte, _ := json.MarshalIndent(response, "", "    ")
			log.Error().Msg(response.Error_message)
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write(response_byte); err != nil {
				log.Error().Msg(errs.ErrWritingToRespBody)
				return
			}
			return
		}
		logMessage += "Регистрация пользователя успешна. Добавление в базу!\n"
		log.Info().Msg(logMessage)
		return // выходим, я больше не отправляю такие сообщения
	}

	var want_email bool
	var want_telegram bool
	var want_webhook []string

	emails_to_send := make([]string, 0)

	// Получаем значения переключателей куда мы хотим получить уведомление
	err, want_email, want_telegram, want_webhook = database_conn_dto.GetCheckboxSettings(req.Notify_Type)
	if err != nil {
		logMessage += fmt.Sprintf("При обращении к базе данных произошла ошибка %v\n", err)
		log.Error().Msg(logMessage)
		return
	}

	notification_message := req.Message

	// Если мы хотим уведомление в ТГ
	if want_telegram {
		if err := tgbot.SendMessage(d.bot, context.Background(), tg_user_id, notification_message); err != nil {
			logMessage += fmt.Sprintf("Ошибка при отправке пользователю уведомления в Telegram: %v", err)
		} else {
			logMessage += fmt.Sprintf("Сообщение было доставлено в Telegram пользователю с id %v\n", tg_user_id)
		}
	}

	// Если мы хотим уведомление по Email
	if want_email {
		emails_to_send = append(emails_to_send, user_email)
		if err := d.smtp.SendMessage(emails_to_send, notification_message, req.Notify_Type); err != nil {
			logMessage += fmt.Sprintf("Ошибка при отправке сообщения на почту: %v\n", err)
		} else {
			logMessage += fmt.Sprintf("Сообщение было доставлено на почту %v\n", user_email)
		}
	}

	// Если мы хотим уведомление по Webhook
	if len(want_webhook) != 0 {
		for _, url := range want_webhook {
			if err := webhook_handler.SendWebhookMessage(url, []byte(notification_message)); err != nil {
				logMessage += fmt.Sprintf("Ошибка при отправке сообщения на адрес %v: %v\n", url, err)
			} else {
				logMessage += fmt.Sprintf("Сообщение было доставлено на Webhook с URL %v\n", url)
			}
		}
	}

	response.Success = true
	response.Error_message = ""
	response_byte, _ := json.MarshalIndent(response, "", "    ")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response_byte); err != nil {
		log.Error().Msg(errs.ErrWritingToRespBody)
		return
	}
	logMessage += "Отправка сообщений была успешно завершена!"
	log.Info().Msg(logMessage)
}
