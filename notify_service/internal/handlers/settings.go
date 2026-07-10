package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"traineesheep/notifyservice/internal/database"
	"traineesheep/notifyservice/internal/errs"
	"traineesheep/notifyservice/internal/types"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// HandleGetNotifySettings обрабатывает GET-запросы на получение текущих настроек уведомлений.
//
// Endpoint: GET /api/get_notify_settings
//
// Функция извлекает userID из JWT токена и возвращает настройки уведомлений пользователя.
//
// Возможные коды ответа:
//   - 200: успешное получение настроек (в формате JSON)
//   - 401: невалидный или отсутствующий JWT токен
//   - 500: внутренняя ошибка сервера

// HandleGetNotifySettings godoc
// @Summary Получение настроек уведомлений
// @Description Возвращает список всех типов уведомлений с их настройками (Telegram, Email, Webhook).
// @Description Требует JWT-токен в заголовке Authorization.
// @Tags Settings
// @Accept json
// @Produce json
// @Success 200 {object} types.NotifyTypeMessengerList "Список настроек"
// @Failure 401 {object} types.ResponseData "Неавторизован"
// @Failure 500 {object} types.ResponseData "Внутренняя ошибка сервера"
// @Router /api/get_notify_settings [get]
// @Security ApiKeyAuth
func HandleGetNotifySettings(w http.ResponseWriter, r *http.Request) {
	var response types.ResponseData
	uid := ulid.Make()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	enableCors(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	logMessage = "Был получен запрос на получение сохраненных настроек админ-панели\n"
	log.Info().Msg(fmt.Sprintf("%v %v", uid.String(), logMessage))

	var jsonList = make([]types.NotifyTypeMessenger, 0)
	err, jsonList := database.GetSettings(types.Ctx, jsonList)
	if err != nil {
		logMessage = fmt.Sprintf("При получении настроек админ-панели из базы произошла ошибка %v", err)
		log.Error().Msg(fmt.Sprintf("%v %v", uid.String(), logMessage))
		return
	}

	var jsonDataList = types.NotifyTypeMessengerList{
		Data: jsonList,
	}

	respByte, err := json.MarshalIndent(jsonDataList, "", "    ")
	if err != nil {
		response.Success = false
		response.ErrorMessage = errs.ErrJsonMarshal + err.Error()
		respBytes, _ := json.MarshalIndent(response, "", "    ")
		logMessage = fmt.Sprint(response.ErrorMessage + "\n")
		log.Error().Msg(fmt.Sprintf("%v %v", uid.String(), logMessage))
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write(respBytes); err != nil {
			logMessage = errs.ErrWritingToRespBody
			log.Error().Msg(fmt.Sprintf("%v %v", uid.String(), logMessage))
			return
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(respByte); err != nil {
		log.Error().Msg(errs.ErrWritingToRespBody)
		return
	}

	logMessage = "Настройки админ-панели были успешно получены из БД!\n"
	log.Info().Msg(fmt.Sprintf("%v %v", uid.String(), logMessage))
}

// HandleSaveSettingsCheckmarks обрабатывает POST-запросы на сохранение настроек уведомлений.
//
// Endpoint: POST /api/notify_types
//
// Тело запроса должно содержать JSON с настройками типов уведомлений для пользователя.
// Функция проверяет JWT токен и сохраняет настройки в базе данных.
//
// Возможные коды ответа:
//   - 200: настройки успешно сохранены
//   - 400: неверный формат запроса
//   - 401: невалидный или отсутствующий JWT токен
//   - 500: внутренняя ошибка сервера

// HandleSaveSettingsCheckmarks godoc
// @Summary Сохранение настроек уведомлений
// @Description Сохраняет настройки типов уведомлений (чекбоксы Telegram, Email, Webhook URL) в БД.
// @Description Требует JWT-токен в заголовке Authorization.
// @Tags Settings
// @Accept json
// @Produce json
// @Param request body types.NotifyTypeMessengerList true "Список настроек для каждого типа уведомления"
// @Success 200 {object} types.ResponseData "Настройки успешно сохранены"
// @Failure 400 {object} types.ResponseData "Неверный формат запроса или ошибка при сохранении"
// @Failure 401 {object} types.ResponseData "Неавторизован (отсутствует или невалидный JWT токен)"
// @Failure 500 {object} types.ResponseData "Внутренняя ошибка сервера"
// @Router /api/notify_types [post]
// @Security ApiKeyAuth
func HandleSaveSettingsCheckmarks(w http.ResponseWriter, r *http.Request) {
	var response types.ResponseData
	uid := ulid.Make()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	enableCors(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	logMessage = "Был получен запрос на сохранение настроек админ-панели в БД\n"
	log.Info().Msg(fmt.Sprintf("%v %v", uid.String(), logMessage))

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		response.Success = false
		response.ErrorMessage = errs.ErrReadingRequestMessage + err.Error()
		respBytes, _ := json.MarshalIndent(response, "", "    ")
		w.WriteHeader(http.StatusInternalServerError)
		logMessage = response.ErrorMessage
		log.Error().Msg(fmt.Sprintf("%v %v", uid.String(), logMessage))
		if _, err := w.Write(respBytes); err != nil {
			logMessage = errs.ErrWritingToRespBody
			log.Error().Msg(fmt.Sprintf("%v %v", uid.String(), logMessage))
			return
		}
		return
	}

	var jsonList types.NotifyTypeMessengerList
	if err := json.Unmarshal(bodyByte, &jsonList); err != nil {
		response.Success = false
		response.ErrorMessage = errs.ErrJsonUnmarshal + err.Error()
		respBytes, _ := json.MarshalIndent(response, "", "    ")
		w.WriteHeader(http.StatusBadRequest)
		log.Error().Msg(fmt.Sprintf("%v %v", uid.String(), response.ErrorMessage))
		if _, err := w.Write(respBytes); err != nil {
			logMessage = errs.ErrWritingToRespBody
			log.Error().Msg(fmt.Sprintf("%v %v", uid.String(), logMessage))
			return
		}
		return
	}

	// После того как получили данные записываем их в базу, не создавая дубликаты
	if err := database.DeleteSettings(types.Ctx); err != nil {
		logMessage = fmt.Sprintf("При обращении к базе данных произошла ошибка: %v", err)
		log.Error().Msg(fmt.Sprintf("%v %v", uid.String(), logMessage))
		return
	}
	for _, elem := range jsonList.Data {

		// Обновляем список разрешенных NotifyType
		notifyTypesAllowed = append(notifyTypesAllowed, elem.NotifyType)

		if err := database.SaveSettings(types.Ctx, elem); err != nil {
			logMessage = fmt.Sprintf("При обращении к базе данных произошла ошибка: %v", err)
			log.Error().Msg(fmt.Sprintf("%v %v", uid.String(), logMessage))
			return
		}
	}

	response.Success = true
	response.ErrorMessage = ""
	respBytes, _ := json.MarshalIndent(response, "", "    ")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(respBytes); err != nil {
		logMessage = errs.ErrWritingToRespBody
		log.Error().Msg(fmt.Sprintf("%v %v", uid.String(), logMessage))
		return
	}

	logMessage = "Настройки админ-панели были успешно сохранены в БД\n"
	log.Info().Msg(fmt.Sprintf("%v %v", uid.String(), logMessage))
}
