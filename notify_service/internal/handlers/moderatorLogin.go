package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"traineesheep/notifyservice/internal/errs"
	"traineesheep/notifyservice/internal/types"

	"github.com/golang-jwt/jwt/v5"
)

// HandleModeratorLogin обрабатывает POST-запросы на авторизацию модератора.
//
// Endpoint: POST /api/moderator_login
//
// Тело запроса должно содержать логин и пароль модератора.
// При успешной авторизации возвращает JWT токен для дальнейших запросов.
//
// Возможные коды ответа:
//   - 200: успешная авторизация (возвращает JWT токен)
//   - 400: неверный формат запроса
//   - 401: неверные учетные данные
//   - 500: внутренняя ошибка сервера

// HandleModeratorLogin godoc
// @Summary Авторизация админа
// @Description Проверяет логин и пароль, возвращает JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body types.LoginData true "Логин и пароль"
// @Success 200 {object} types.ResponseData
// @Failure 400 {object} types.ResponseData
// @Failure 401 {object} types.ResponseData
// @Router /api/moderator_login [post]
func (d DTO) HandleModeratorLogin(w http.ResponseWriter, r *http.Request) {
	var response types.ResponseData
	var login_data types.LoginData
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

	logMessage += "Был получен запрос на вход пользователя в админ-панель\n"
	logMessage += "=========================\n"

	body_byte, err := io.ReadAll(r.Body)
	if err != nil {
		response.Success = false
		response.Error_message = errs.ErrReadingRequestMessage + err.Error()
		response_byte, _ := json.MarshalIndent(response, "", "    ")
		logMessage += response.Error_message + "\n"
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write(response_byte); err != nil {
			log.Println(errs.ErrWritingToRespBody)
		}
		return
	}

	if err := json.Unmarshal(body_byte, &login_data); err != nil {

		var response types.ResponseData
		response.Success = false
		response.Error_message = errs.ErrJsonUnmarshal + ": " + err.Error()
		response_byte, _ := json.MarshalIndent(response, "", "    ")
		logMessage += response.Error_message + "\n"
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write(response_byte); err != nil {
			log.Println(errs.ErrWritingToRespBody)
		}
		return
	}

	if login_data.Login == d.admin_login && login_data.Password == d.admin_pass {

		jwt_key = []byte(d.JwtSecret)
		jwt_t = jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{
				"iss": "my_auth_server",
				"sub": "admin",
				"exp": time.Now().Add(15 * time.Minute).Unix(),
			})
		jwt_s, err := jwt_t.SignedString(jwt_key)
		if err != nil {
			logMessage += fmt.Sprintf("Ошибка при создании криптографической подписи для токена: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var response types.ResponseData
		response.Success = true
		response.Error_message = ""
		response.Token = jwt_s
		response_byte, _ := json.MarshalIndent(response, "", "    ")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(response_byte); err != nil {
			log.Println(errs.ErrWritingToRespBody)
		}
		logMessage += "Логин и пароль успешно прошли проверку!\n"
		log.Println(logMessage)
		return
	} else {
		var response types.ResponseData
		response.Success = false
		response.Error_message = "Неправильно введен логин или пароль. Попробуйте еще раз!"
		response_byte, _ := json.MarshalIndent(response, "", "    ")
		logMessage += response.Error_message + "\n"
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(response_byte); err != nil {
			log.Println(errs.ErrWritingToRespBody)
		}
		return
	}
}
