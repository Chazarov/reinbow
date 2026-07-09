package interfaces

import (
	"context"

	"github.com/tokyobordel/traineepkg/errors"
)

// INotifyService описывает отправку уведомлений о модерации.
type INotifyService interface {
	NotifyImageModeration(ctx context.Context, imageID int) errors.DomainError
}
