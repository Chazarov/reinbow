package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"traineesheep/notifyservice/internal/database"
	"traineesheep/notifyservice/internal/errs"
	"traineesheep/notifyservice/internal/types"
)

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
func (d DTO) HandleSaveSettingsCheckmarks(w http.ResponseWriter, r *http.Request) {
	var response types.ResponseData
	database_conn_dto := database.NewDatabaseDTO(d.sql_connection)

	enableCors(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	logMessage += "Был получен запрос на сохранение настроек админ-панели в БД\n"
	logMessage += "=========================\n"

	// _, err := ValidateToken(r, d.JwtSecret)
	// if err != nil {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	response.Success = false
	// 	response.Error_message = "Unauthorized"
	// 	respBytes, _ := json.MarshalIndent(response, "", "    ")
	// 	logMessage += "Пользователь был не авторизован. Действие не выполнено!\n"
	// 	if _, err := w.Write(respBytes); err != nil {
	// 		log.Println("Failed to write response:", err)
	// 		return
	// 	}
	// 	return
	// }

	body_byte, err := io.ReadAll(r.Body)
	if err != nil {
		response.Success = false
		response.Error_message = errs.ErrReadingRequestMessage + err.Error()
		respBytes, _ := json.MarshalIndent(response, "", "    ")
		logMessage += response.Error_message + "\n"
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write(respBytes); err != nil {
			log.Println(errs.ErrWritingToRespBody)
			return
		}
		return
	}

	var json_list types.NotifyTypeMessengerList
	if err := json.Unmarshal(body_byte, &json_list); err != nil {
		response.Success = false
		response.Error_message = errs.ErrJsonUnmarshal + err.Error()
		respBytes, _ := json.MarshalIndent(response, "", "    ")
		logMessage += response.Error_message + "\n"
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		if _, err := w.Write(respBytes); err != nil {
			log.Println(errs.ErrWritingToRespBody)
			return
		}
		return
	}

	fmt.Println("Мне пришло", json_list)

	// После того как получили данные записываем их в базу, не создавая дубликаты
	if err := database_conn_dto.DeleteSettings(); err != nil {
		logMessage += fmt.Sprintf("При обращении к базе данных произошла ошибка: %v", err)
		return
	}
	for _, elem := range json_list.Data {

		// Обновляем список разрешенных NotifyType
		notify_types_allowed = append(notify_types_allowed, elem.NotifyType)

		if err := database_conn_dto.SaveSettings(elem); err != nil {
			logMessage += fmt.Sprintf("При обращении к базе данных произошла ошибка: %v", err)
			return
		}
	}

	response.Success = true
	response.Error_message = ""
	respBytes, _ := json.MarshalIndent(response, "", "    ")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(respBytes); err != nil {
		log.Println(errs.ErrWritingToRespBody)
		return
	}

	logMessage += "Настройки админ-панели были успешно сохранены в БД\n"
	log.Println(logMessage)
}
