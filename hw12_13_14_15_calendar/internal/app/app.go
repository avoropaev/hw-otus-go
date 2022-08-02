package app

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrEventNotFound           = errors.New("event not found")
	ErrNotificationAlreadySent = errors.New("notification already sent")
	RemoveOlderThan            = time.Hour * 24 * 365
)

type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, eventGUID uuid.UUID, event storage.Event) error
	DeleteEvent(ctx context.Context, eventGUID uuid.UUID) error
	GetEventForDay(ctx context.Context, startDateTime time.Time) ([]*storage.Event, error)
	GetEventForWeek(ctx context.Context, startDateTime time.Time) ([]*storage.Event, error)
	GetEventForMonth(ctx context.Context, startDateTime time.Time) ([]*storage.Event, error)
	GetEventForNotify(ctx context.Context) ([]*storage.Event, error)
	SendNotificationAndMarkAsNotifier(ctx context.Context, guid uuid.UUID, title string, startAt time.Time, userGUID uuid.UUID) error
	RemoveOldEvents(ctx context.Context) (int64, error)
}

type app struct {
	storage  Storage
	notifier Notifier
}

type Storage interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, eventGUID uuid.UUID, event storage.Event) error
	DeleteEvent(ctx context.Context, eventGUID uuid.UUID) error
	FindEventsByInterval(ctx context.Context, startDateTime, endDateTime time.Time) ([]*storage.Event, error)
	FindEventByGUID(ctx context.Context, eventGUID uuid.UUID) (*storage.Event, error)
	FindEventsNeedsNotify(ctx context.Context) ([]*storage.Event, error)
	DeleteEventsOlderThan(ctx context.Context, datetime time.Time) (rowAffected int64, err error)
}

func New(storage Storage) Application {
	return &app{
		storage:  storage,
		notifier: NewNotifier(),
	}
}

func (a *app) CreateEvent(ctx context.Context, event storage.Event) error {
	return a.storage.CreateEvent(ctx, event)
}

func (a *app) UpdateEvent(ctx context.Context, eventGUID uuid.UUID, event storage.Event) error {
	return a.storage.UpdateEvent(ctx, eventGUID, event)
}

func (a *app) DeleteEvent(ctx context.Context, eventGUID uuid.UUID) error {
	return a.storage.DeleteEvent(ctx, eventGUID)
}

func (a *app) GetEventForDay(ctx context.Context, startDateTime time.Time) ([]*storage.Event, error) {
	return a.storage.FindEventsByInterval(ctx, startDateTime, startDateTime.Add(time.Hour*24))
}

func (a *app) GetEventForWeek(ctx context.Context, startDateTime time.Time) ([]*storage.Event, error) {
	return a.storage.FindEventsByInterval(ctx, startDateTime, startDateTime.Add(time.Hour*24*7))
}

func (a *app) GetEventForMonth(ctx context.Context, startDateTime time.Time) ([]*storage.Event, error) {
	return a.storage.FindEventsByInterval(ctx, startDateTime, startDateTime.Add(time.Hour*24*30))
}

func (a *app) GetEventForNotify(ctx context.Context) ([]*storage.Event, error) {
	return a.storage.FindEventsNeedsNotify(ctx)
}

func (a *app) SendNotificationAndMarkAsNotifier(
	ctx context.Context,
	eventGUID uuid.UUID,
	title string,
	startAt time.Time,
	userGUID uuid.UUID,
) error {
	event, err := a.storage.FindEventByGUID(ctx, eventGUID)
	if err != nil {
		return err
	}

	if event.Notified {
		return ErrNotificationAlreadySent
	}

	if event == nil {
		return ErrEventNotFound
	}

	a.notifier.Notify(ctx, eventGUID, title, startAt, userGUID)

	event.Notified = true

	err = a.storage.UpdateEvent(ctx, eventGUID, *event)
	if err != nil {
		return err
	}

	return nil
}

func (a *app) RemoveOldEvents(ctx context.Context) (int64, error) {
	return a.storage.DeleteEventsOlderThan(ctx, time.Now().Add(RemoveOlderThan*-1))
}
