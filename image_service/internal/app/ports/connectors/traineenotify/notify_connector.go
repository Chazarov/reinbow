package traineenotify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tokyobordel/traineepkg/errors"
	"github.com/tokyobordel/traineepkg/logger"
)

// NotifyConnector отправляет уведомления через внешний сервис.
type NotifyConnector struct {
	notifyServiceURL string
	externalURL      string
	logger           *logger.ContextLogger
	client           *http.Client
	tgAdminId        int
}

// NewNotificatorService создаёт клиент внешнего сервиса уведомлений.
func NewNotificatorService(notifyServiceURL, externalURL string, tgAdminId int, logger *logger.ContextLogger) (*NotifyConnector, error) {
	if strings.TrimSpace(notifyServiceURL) == "" {
		return nil, fmt.Errorf("notify service url is required")
	}
	if strings.TrimSpace(externalURL) == "" {
		return nil, fmt.Errorf("external url is required")
	}

	return &NotifyConnector{
		notifyServiceURL: strings.TrimRight(notifyServiceURL, "/"),
		externalURL:      strings.TrimRight(externalURL, "/"),
		logger:           logger,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		tgAdminId: tgAdminId,
	}, nil
}

// moderationLink формирует ссылку на страницу модерации.
func (c *NotifyConnector) moderationLink(imageID int) string {
	return fmt.Sprintf("%s", c.externalURL)
}

// NotifyImageModeration отправляет уведомление о необходимости модерации изображения.
func (c *NotifyConnector) NotifyImageModeration(ctx context.Context, imageID int) errors.DomainError {
	message := fmt.Sprintf(
		"Новое изображение требует модерации\nimage_id: %d\nперейти: %s",
		imageID,
		c.moderationLink(imageID),
	)
	return c.sendNotify(ctx, message, AdminNewImg)
}

// sendNotify отправляет HTTP-запрос во внешний сервис уведомлений.
func (c *NotifyConnector) sendNotify(ctx context.Context, message string, NType NotifyType) errors.DomainError {
	reqBody := NotifyRequest{
		Message:    message,
		NotifyType: NType,
		TelegramId: c.tgAdminId,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		c.logger.Criticalf(ctx, "Failed to marshal notify request: %v", err)
		return errors.NewInternalServiceError("failed to marshal notification request", err)
	}

	url := fmt.Sprintf("%s/api/notify", c.notifyServiceURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Criticalf(ctx, "Failed to create notify request: %v", err)
		return errors.NewInternalServiceError("failed to create notification request", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Criticalf(ctx, "Failed to send notify request: %v", err)
		return errors.NewInternalServiceError("failed to send notification request", err)
	}
	defer resp.Body.Close()

	var notifyResp NotifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&notifyResp); err != nil {
		c.logger.Criticalf(ctx, "Failed to decode notify response: %v", err)
		return errors.NewInternalServiceError("failed to decode notification response", err)
	}

	if !notifyResp.Success {
		c.logger.Criticalf(ctx, "Notification service returned error: %s", notifyResp.ErrMessage)
		return errors.NewIntegrationError(fmt.Sprintf("notification failed: %s", notifyResp.ErrMessage), nil)
	}

	c.logger.Infof(ctx, "Notification sent successfully: %s", message)
	return nil
}
