package database

import (
	"context"
	"traineesheep/notifyservice/internal/types"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseDTO struct {
	sql_connection *pgxpool.Pool
}

func NewDatabaseDTO(conn *pgxpool.Pool) *DatabaseDTO {
	return &DatabaseDTO{sql_connection: conn}
}

type CheckboxesDTO struct {
	want_email    bool
	want_telegram bool
	want_webhook  bool
}

func (d DatabaseDTO) AddEmail(email string) error {
	conn := d.sql_connection
	query := `INSERT INTO client (email) VALUES ($1)`
	_, err := conn.Exec(context.Background(), query, email)
	if err != nil {
		return err
	}
	return nil
}

func (d DatabaseDTO) GetCheckboxSettings(notify_type string) (error, bool, bool, []string) {
	var want_email, want_telegram bool
	var want_webhook []string
	conn := d.sql_connection
	query := `SELECT want_email, want_telegram, want_webhook FROM notify_type_message
	WHERE notify_type = $1`
	row := conn.QueryRow(context.Background(), query, notify_type)
	if err := row.Scan(&want_email, &want_telegram, &want_webhook); err != nil {
		return err, false, false, nil
	}
	return nil, want_email, want_telegram, want_webhook
}

func (d DatabaseDTO) GetSettings(json_list []types.NotifyTypeMessenger) (error, []types.NotifyTypeMessenger) {
	conn := d.sql_connection

	query := `SELECT notify_type, notify_description, want_telegram, want_email, want_webhook FROM notify_type_message`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	for rows.Next() {
		var notify_type string
		var notify_description string
		var want_email bool
		var want_telegram bool
		var webhook_urls []string
		if err := rows.Scan(&notify_type, &notify_description, &want_telegram, &want_email, &webhook_urls); err != nil {
			return err, nil
		}
		json_list = append(json_list, types.NotifyTypeMessenger{
			NotifyType:   notify_type,
			Description:  notify_description,
			WantEmail:    want_email,
			WantTelegram: want_telegram,
			WebhookUrls:  webhook_urls,
		})
	}
	return nil, json_list
}

func (d DatabaseDTO) SaveSettings(elem types.NotifyTypeMessenger) error {
	conn := d.sql_connection

	query := `INSERT INTO notify_type_message (notify_type, notify_description, want_telegram, want_email, want_webhook)
    		VALUES ($1, $2, $3, $4, $5)
    		ON CONFLICT (notify_type) DO UPDATE SET
			notify_description = EXCLUDED.notify_description,
        	want_telegram = EXCLUDED.want_telegram,
        	want_email = EXCLUDED.want_email,
			want_webhook = EXCLUDED.want_webhook;`

	_, err := conn.Exec(context.Background(), query, elem.NotifyType, elem.Description, elem.WantTelegram, elem.WantEmail, elem.WebhookUrls)
	if err != nil {
		return err
	}
	return nil
}

func (d DatabaseDTO) DeleteSettings() error {
	conn := d.sql_connection
	query := `TRUNCATE notify_type_message`
	if _, err := conn.Exec(context.Background(), query); err != nil {
		return err
	}
	return nil
}
