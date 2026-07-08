package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"traineesheep/notifyservice/internal/database"
	"traineesheep/notifyservice/internal/errs"
	"traineesheep/notifyservice/internal/types"
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
func (d DTO) HandleGetNotifySettings(w http.ResponseWriter, r *http.Request) {
	var response types.ResponseData
	database_conn_dto := database.NewDatabaseDTO(d.sql_connection)

	enableCors(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	logMessage += "Был получен запрос на получение сохраненных настроек админ-панели\n"
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

	var json_list = make([]types.NotifyTypeMessenger, 0)
	err, json_list := database_conn_dto.GetSettings(json_list)
	if err != nil {
		logMessage += fmt.Sprintf("При получении настроек админ-панели из базы произошла ошибка %v", err)
		return
	}

	var json_data_list = types.NotifyTypeMessengerList{
		Data: json_list,
	}

	resp_byte, err := json.MarshalIndent(json_data_list, "", "    ")
	if err != nil {
		response.Success = false
		response.Error_message = errs.ErrJsonMarshal + err.Error()
		respBytes, _ := json.MarshalIndent(response, "", "    ")
		logMessage += fmt.Sprintf(response.Error_message + "\n")
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write(respBytes); err != nil {
			log.Println(errs.ErrWritingToRespBody)
			return
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp_byte); err != nil {
		log.Println(errs.ErrWritingToRespBody)
		return
	}

	logMessage += "Настройки админ-панели были успешно получены из БД!\n"
	log.Println(logMessage)
}
