package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"traineesheep/notifyservice/internal/database"
	"traineesheep/notifyservice/internal/errs"
	"traineesheep/notifyservice/internal/tgbot"
	"traineesheep/notifyservice/internal/types"
	"traineesheep/notifyservice/internal/webhook_handler"
	"traineesheep/notifyservice/pkg/email"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logMessage string

var TelegramID int64

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
func HandleNotify(w http.ResponseWriter, r *http.Request) {
	uid := ulid.Make()
	var response types.ResponseData
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	var req types.Recipent

	enableCors(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// --- Чтение и валидация запроса (синхронно, до захвата слота) ---
	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		response.Success = false
		response.ErrorMessage = errs.ErrReadingRequestMessage + err.Error()
		respBytes, _ := json.MarshalIndent(response, "", "    ")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBytes)
		return
	}

	if err := json.Unmarshal(bodyData, &req); err != nil {
		response.Success = false
		response.ErrorMessage = errs.ErrJsonUnmarshal + err.Error()
		respBytes, _ := json.MarshalIndent(response, "", "    ")
		log.Error().Msg(fmt.Sprintf("%v %v", uid.String(), response.ErrorMessage))
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	log.Info().Msg(fmt.Sprintf("%v Получен запрос на тип %v", uid.String(), req.NotifyType))
	log.Info().Msg(fmt.Sprintf("%v Сообщение: \"%v\"", uid.String(), req.Message))

	// --- Проверка допустимости типа уведомления ---
	isAllowed := false
	rows, err := database.GetNotifyTypes(types.Ctx)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("%v %v", uid.String(), err.Error()))
		response.Success = false
		response.ErrorMessage = "Ошибка получения типов уведомлений"
		respBytes, _ := json.MarshalIndent(response, "", "    ")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBytes)
		return
	}
	if len(notifyTypesAllowed) == 0 {
		for rows.Next() {
			var temp string
			rows.Scan(&temp)
			notifyTypesAllowed = append(notifyTypesAllowed, temp)
		}
	}
	for _, elem := range notifyTypesAllowed {
		if elem == req.NotifyType {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		response.Success = false
		response.ErrorMessage = fmt.Sprintf("Недопустимый тип уведомления! %v", req.NotifyType)
		respBytes, _ := json.MarshalIndent(response, "", "    ")
		log.Error().Msg(fmt.Sprintf("%v %v", uid.String(), response.ErrorMessage))
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	// --- Обработка регистрации (выполняется синхронно, без ограничения) ---
	if req.NotifyType == "user_register" {
		if err := database.AddEmail(types.Ctx, req.Email); err != nil {
			response.Success = false
			response.ErrorMessage = "При попытке добавить данные в базу произошла ошибка: " + err.Error()
			respBytes, _ := json.MarshalIndent(response, "", "    ")
			log.Error().Msg(fmt.Sprintf("%v %v", uid.String(), response.ErrorMessage))
			w.WriteHeader(http.StatusBadRequest)
			w.Write(respBytes)
			return
		}
		response.Success = true
		response.ErrorMessage = "Пользователь зарегистрирован"
		respBytes, _ := json.MarshalIndent(response, "", "    ")
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
		log.Info().Msg(fmt.Sprintf("%v Регистрация пользователя успешна", uid.String()))
		return
	}

	// --- Основная задача: отправка уведомлений (ограничивается через errgroup) ---
	// Используем WaitGroup для синхронизации внутри хендлера
	var wg sync.WaitGroup
	wg.Add(1)

	// Запускаем задачу через errgroup.Go – блокируется, если достигнут лимит
	TaskGroup.Go(func() error {
		defer wg.Done()
		// Логика отправки (копия из исходного кода)
		var checks types.CheckboxesParams
		emailsToSend := make([]string, 0)

		err, checks := database.GetCheckboxSettings(types.Ctx, req.NotifyType)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("%v Ошибка получения настроек: %v", uid.String(), err))
			return nil // не возвращаем ошибку, чтобы не останавливать группу
		}

		if checks.WantTelegram {
			if err := tgbot.SendMessage(TgBot, types.Ctx, TelegramID, req.Message); err != nil {
				log.Error().Msg(fmt.Sprintf("%v Ошибка отправки в Telegram: %v", uid.String(), err))
			} else {
				log.Info().Msg(fmt.Sprintf("%v Сообщение доставлено в Telegram", uid.String()))
			}
		}

		if checks.WantEmail {
			emailsToSend = append(emailsToSend, req.Email)
			if err := email.SendMessage(emailsToSend, req.Message, req.NotifyType); err != nil {
				log.Error().Msg(fmt.Sprintf("%v Ошибка отправки email: %v", uid.String(), err))
			} else {
				log.Info().Msg(fmt.Sprintf("%v Сообщение доставлено на почту %v", uid.String(), req.Email))
			}
		}

		if len(checks.WantWebhook) != 0 {
			for _, url := range checks.WantWebhook {
				if err := webhook_handler.SendWebhookMessage(url, []byte(req.Message)); err != nil {
					log.Error().Msg(fmt.Sprintf("%v Ошибка отправки webhook %v: %v", uid.String(), url, err))
				} else {
					log.Info().Msg(fmt.Sprintf("%v Сообщение доставлено на webhook %v", uid.String(), url))
				}
			}
		}

		log.Info().Msg(fmt.Sprintf("%v Отправка сообщений завершена", uid.String()))
		return nil
	})

	// Ожидаем завершения задачи (синхронизация)
	wg.Wait()

	// Отправляем успешный ответ после выполнения
	response.Success = true
	response.ErrorMessage = ""
	respBytes, _ := json.MarshalIndent(response, "", "    ")
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
	log.Info().Msg(fmt.Sprintf("%v Запрос обработан успешно", uid.String()))
}
