package app

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	storage Storage
}

type Storage interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, eventGUID uuid.UUID, event storage.Event) error
	DeleteEvent(ctx context.Context, eventGUID uuid.UUID) error
	FindEventsByInterval(ctx context.Context, startDateTime, endDateTime time.Time) ([]*storage.Event, error)
	FindEventByGUID(ctx context.Context, eventGUID uuid.UUID) (*storage.Event, error)
}

func New(storage Storage) *App {
	return &App{storage}
}

func (a *App) TestDB(ctx context.Context) error {
	notifyBefore := time.Hour

	eventGUID := uuid.New()

	err := a.storage.CreateEvent(ctx, storage.Event{
		GUID:         eventGUID,
		Title:        "title",
		StartAt:      time.Now(),
		EndAt:        time.Now().Add(time.Hour * 5),
		UserGUID:     uuid.MustParse("1d8fe576-d420-479b-b96f-30fd7e0107c1"),
		NotifyBefore: &notifyBefore,
	})
	if err != nil {
		return err
	}

	err = a.storage.UpdateEvent(ctx, eventGUID, storage.Event{
		Title:    "new title",
		StartAt:  time.Now().Add(time.Hour * 5),
		EndAt:    time.Now().Add(time.Hour * 10),
		UserGUID: uuid.MustParse("1d8fe576-d420-479b-b96f-30fd7e0107c1"),
	})
	if err != nil {
		return err
	}

	events, err := a.storage.FindEventsByInterval(ctx, time.Now(), time.Now().Add(time.Hour*24))
	if err != nil {
		return err
	}

	_ = events

	err = a.storage.DeleteEvent(ctx, eventGUID)
	if err != nil {
		return err
	}

	events, err = a.storage.FindEventsByInterval(ctx, time.Now(), time.Now().Add(time.Hour*24))
	if err != nil {
		return err
	}

	_ = events

	return nil
}
