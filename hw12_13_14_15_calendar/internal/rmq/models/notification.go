package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	EventGUID     uuid.UUID `json:"event_guid"`
	EventTitle    string    `json:"event_title"`
	EventStartAt  time.Time `json:"event_start_at"`
	EventUserGUID uuid.UUID `json:"event_user_guid"`
}
