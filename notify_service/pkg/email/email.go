package email

// Пакет email используется для реализации логики отправки уведомлений пользователям
// на их почту(-ы)

import (
	"crypto/tls"
	"log"
	"net/smtp"
	"strings"
)

// Структура SmtpDTO нужна для передачи параметров почты, с которой будет
// вестись рассылка
type SmtpDTO struct {
	Email    string // Почта для рассылки
	Password string // Пароль от почты
	Host     string // Хост
	Port     string // Порт
}

// Функция NewSmtpDTO используется для создания экземпляра структуры NewSmtpDTO
// Она возращает созданный экземпляр стурктуры NewSmtpDTO
func NewSmtpDTO(e string, p string, h string, port string) *SmtpDTO {
	return &SmtpDTO{
		Email:    e,
		Password: p,
		Host:     h,
		Port:     port,
	}
}

// Функция SendMessage используется для отправки писем с уведомлениями
// На вход получаем:
// receiverEmails - массив из адресов почт, куда придут уведмоения
// message - текст письма
// notify_type - тип уведомления
// На выходе получаем лог об успешности отправки
func (s SmtpDTO) SendMessage(receiverEmails []string, message string, notify_type string) error {
	tlsConfig := &tls.Config{
		ServerName: s.Host,
	}

	conn, err := tls.Dial("tcp", s.Host+":"+s.Port, tlsConfig)
	if err != nil {
		//*logMessage += fmt.Sprintf("Ошибка при подключении TLS: %v", err.Error())
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.Host)
	if err != nil {
		//*logMessage += fmt.Sprintf("Ошибка при создании клиента: %v", err)
		return err
	}
	defer client.Quit()

	auth := smtp.PlainAuth("", s.Email, s.Password, s.Host)
	if err = client.Auth(auth); err != nil {
		//log.Printf("аутентификация: %v", err)
		return err
	}

	if err = client.Mail(s.Email); err != nil {
		//log.Printf("отправитель: %v", err)
		return err
	}

	for _, rcpt := range receiverEmails {
		if err = client.Rcpt(rcpt); err != nil {
			//log.Printf("получатель %s: %v", rcpt, err)
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		//log.Printf("открытие Data: %v", err)
		return err
	}
	defer w.Close()

	var topic string
	switch notify_type {
	case "user_register":
		topic = "Регистрация аккаунта"
	case "user_login":
		topic = "Вход в аккаунт"
	case "admin_newImg":
		topic = "Новое изображение для модерации"
	case "user_imgVerdict":
		topic = "Модерация вашего поста"
	case "user_confirm":
		topic = "Подтверждение почты"
	default:
		topic = "Служебное сообщение"
	}

	var temp_message = "Subject: " + topic + "\r\n" + "\r\n"
	temp_message += message
	temp_message = strings.ReplaceAll(temp_message, "<b>", "")
	temp_message = strings.ReplaceAll(temp_message, "</b>", "")
	byte_temp_message := []byte(temp_message)

	_, err = w.Write(byte_temp_message)
	if err != nil {
		log.Printf("запись письма: %v", err)
		return err
	}

	return nil
}
