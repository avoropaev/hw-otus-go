package app

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, eventGUID uuid.UUID, event storage.Event) error
	DeleteEvent(ctx context.Context, eventGUID uuid.UUID) error
	GetEventForDay(ctx context.Context, startDateTime time.Time) ([]*storage.Event, error)
	GetEventForWeek(ctx context.Context, startDateTime time.Time) ([]*storage.Event, error)
	GetEventForMonth(ctx context.Context, startDateTime time.Time) ([]*storage.Event, error)
}

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

func (a *App) CreateEvent(ctx context.Context, event storage.Event) error {
	return a.storage.CreateEvent(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, eventGUID uuid.UUID, event storage.Event) error {
	return a.storage.UpdateEvent(ctx, eventGUID, event)
}

func (a *App) DeleteEvent(ctx context.Context, eventGUID uuid.UUID) error {
	return a.storage.DeleteEvent(ctx, eventGUID)
}

func (a *App) GetEventForDay(ctx context.Context, startDateTime time.Time) ([]*storage.Event, error) {
	return a.storage.FindEventsByInterval(ctx, startDateTime, startDateTime.Add(time.Hour*24))
}

func (a *App) GetEventForWeek(ctx context.Context, startDateTime time.Time) ([]*storage.Event, error) {
	return a.storage.FindEventsByInterval(ctx, startDateTime, startDateTime.Add(time.Hour*24*7))
}

func (a *App) GetEventForMonth(ctx context.Context, startDateTime time.Time) ([]*storage.Event, error) {
	return a.storage.FindEventsByInterval(ctx, startDateTime, startDateTime.Add(time.Hour*24*30))
}
