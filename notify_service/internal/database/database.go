package database

import (
	"context"
	"traineesheep/notifyservice/internal/types"

	"github.com/jackc/pgx/v5"
)

// Функция AddEmail нужна для добавления адреса почты в базу данных
func AddEmail(Ctx context.Context, email string) error {
	conn := types.SqlConnection
	query := `INSERT INTO client (email) VALUES ($1)`
	_, err := conn.Exec(Ctx, query, email)
	if err != nil {
		return err
	}
	return nil
}

// Функция GetCheckboxSettings нужна для сохранения настроек переключателей из админ-панели в базу
func GetCheckboxSettings(Ctx context.Context, notify_type string) (error, types.CheckboxesParams) {
	var wantEmail, wantTelegram bool
	var wantWebhook []string
	conn := types.SqlConnection
	query := `SELECT want_email, want_telegram, want_webhook FROM notify_type_message
	WHERE notify_type = $1`
	row := conn.QueryRow(Ctx, query, notify_type)
	if err := row.Scan(&wantEmail, &wantTelegram, &wantWebhook); err != nil {
		return err, types.CheckboxesParams{WantEmail: false, WantTelegram: false, WantWebhook: nil}
	}
	return nil, types.CheckboxesParams{WantEmail: wantEmail, WantTelegram: wantTelegram, WantWebhook: wantWebhook}
}

// Функция GetSettings нужна для получения значений таблицы админ-панели из базы
func GetSettings(Ctx context.Context, jsonList []types.NotifyTypeMessenger) (error, []types.NotifyTypeMessenger) {
	conn := types.SqlConnection

	query := `SELECT notify_type, notify_description, want_telegram, want_email, want_webhook FROM notify_type_message`
	rows, err := conn.Query(Ctx, query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	for rows.Next() {
		var notifyType string
		var notifyDescription string
		var wantEmail bool
		var wantTelegram bool
		var webhookUrls []string
		if err := rows.Scan(&notifyType, &notifyDescription, &wantTelegram, &wantEmail, &webhookUrls); err != nil {
			return err, nil
		}
		jsonList = append(jsonList, types.NotifyTypeMessenger{
			NotifyType:   notifyType,
			Description:  notifyDescription,
			WantEmail:    wantEmail,
			WantTelegram: wantTelegram,
			WebhookUrls:  webhookUrls,
		})
	}
	return nil, jsonList
}

// Функция SaveSettings нужна для сохранения отредактированной таблицы админ-панели в базу
func SaveSettings(Ctx context.Context, elem types.NotifyTypeMessenger) error {
	conn := types.SqlConnection

	query := `INSERT INTO notify_type_message (notify_type, notify_description, want_telegram, want_email, want_webhook)
    		VALUES ($1, $2, $3, $4, $5)
    		ON CONFLICT (notify_type) DO UPDATE SET
			notify_description = EXCLUDED.notify_description,
        	want_telegram = EXCLUDED.want_telegram,
        	want_email = EXCLUDED.want_email,
			want_webhook = EXCLUDED.want_webhook;`

	_, err := conn.Exec(Ctx, query, elem.NotifyType, elem.Description, elem.WantTelegram, elem.WantEmail, elem.WebhookUrls)
	if err != nil {
		return err
	}
	return nil
}

// Функция DeleteSettings нужна для очистки таблицы с типами уведомлений в базе
func DeleteSettings(Ctx context.Context) error {
	conn := types.SqlConnection
	query := `TRUNCATE notify_type_message`
	if _, err := conn.Exec(Ctx, query); err != nil {
		return err
	}
	return nil
}

// Функция GetNotifyTypes нужна для выбора всех типов уведомлений из базы
func GetNotifyTypes(Ctx context.Context) (pgx.Rows, error) {
	conn := types.SqlConnection
	query := `SELECT notify_type FROM notify_type_message`
	rows, err := conn.Query(Ctx, query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
