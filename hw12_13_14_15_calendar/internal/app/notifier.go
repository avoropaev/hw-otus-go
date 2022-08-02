package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Notifier interface {
	Notify(ctx context.Context, guid uuid.UUID, title string, startAt time.Time, userGUID uuid.UUID)
}

type notifier struct{}

func NewNotifier() Notifier {
	return &notifier{}
}

func (n *notifier) Notify(_ context.Context, eventGUID uuid.UUID, title string, startAt time.Time, userGUID uuid.UUID) {
	l := log.With().Fields(map[string]interface{}{
		"event_guid": eventGUID.String(),
		"title":      title,
		"start_at":   startAt,
		"userGUID":   userGUID,
	}).Logger()

	l.Info().Msg("Отправляем уведомление...")
}
