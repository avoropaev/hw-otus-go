package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events sync.Map
}

var _ app.Storage = (*Storage)(nil)

func New() *Storage {
	return &Storage{}
}

func (s *Storage) CreateEvent(_ context.Context, event storage.Event) error {
	if _, loaded := s.events.LoadOrStore(event.GUID, event); loaded {
		return storage.ErrEventAlreadyExists
	}

	return nil
}

func (s *Storage) UpdateEvent(_ context.Context, eventGUID uuid.UUID, newEvent storage.Event) error {
	newEvent.GUID = eventGUID

	val, ok := s.events.Load(eventGUID)
	if !ok {
		return storage.ErrEventNotFound
	}

	oldEvent := val.(storage.Event)

	if newEvent != oldEvent {
		s.events.Delete(eventGUID)
		s.events.Store(eventGUID, newEvent)
	}

	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, eventGUID uuid.UUID) error {
	s.events.Delete(eventGUID)

	return nil
}

func (s *Storage) FindEventsByInterval(_ context.Context, start, end time.Time) (events []*storage.Event, err error) {
	s.events.Range(func(key, value interface{}) bool {
		event := value.(storage.Event)

		if event.EndAt.Before(start) {
			return true
		}

		if event.StartAt.After(end) {
			return true
		}

		events = append(events, &event)

		return true
	})

	return events, nil
}

func (s *Storage) FindEventByGUID(_ context.Context, eventGUID uuid.UUID) (*storage.Event, error) {
	if value, ok := s.events.Load(eventGUID); ok {
		event := value.(storage.Event)

		return &event, nil
	}

	return nil, nil
}

func (s *Storage) FindEventsNeedsNotify(_ context.Context) (events []*storage.Event, err error) {
	s.events.Range(func(key, value interface{}) bool {
		event := value.(storage.Event)

		if event.NotifyBefore == nil || event.Notified {
			return true
		}

		notifyAt := event.StartAt.Add(*event.NotifyBefore * -1)
		if notifyAt.After(time.Now()) {
			return true
		}

		events = append(events, &event)

		return true
	})

	return events, nil
}

func (s *Storage) DeleteEventsOlderThan(_ context.Context, datetime time.Time) (rowAffected int64, err error) {
	eventsToDelete := make([]*storage.Event, 0)

	s.events.Range(func(key, value interface{}) bool {
		event := value.(storage.Event)

		if datetime.After(event.StartAt) {
			eventsToDelete = append(eventsToDelete, &event)
		}

		return true
	})

	for _, event := range eventsToDelete {
		s.events.Delete(event.GUID)
	}

	return int64(len(eventsToDelete)), nil
}
