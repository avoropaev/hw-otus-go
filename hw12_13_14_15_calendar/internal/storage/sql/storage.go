package sqlstorage

import (
	"context"
	"errors"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	conn *pgxpool.Pool
}

var _ app.Storage = (*Storage)(nil)

func New(conn *pgxpool.Pool) *Storage {
	return &Storage{conn}
}

func (s Storage) CreateEvent(ctx context.Context, e storage.Event) error {
	event, err := s.FindEventByGUID(ctx, e.GUID)
	if err != nil {
		return err
	}

	if event != nil {
		return storage.ErrEventAlreadyExists
	}

	sql := `
		INSERT INTO events(guid, title, start_at, end_at, description, user_guid, notify_before)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`
	_, err = s.conn.Exec(ctx, sql, e.GUID, e.Title, e.StartAt, e.EndAt, e.Description, e.UserGUID, e.NotifyBefore)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) UpdateEvent(ctx context.Context, eventGUID uuid.UUID, e storage.Event) error {
	event, err := s.FindEventByGUID(ctx, e.GUID)
	if err != nil {
		return err
	}

	if event == nil {
		return storage.ErrEventNotFound
	}

	sql := `
		UPDATE events
		SET title = $2, start_at = $3, end_at = $4, description = $5, user_guid = $6, notify_before = $7, notified = $8
		WHERE guid = $1
	`

	_, err = s.conn.Exec(ctx, sql, eventGUID, e.Title, e.StartAt, e.EndAt, e.Description, e.UserGUID, e.NotifyBefore, e.Notified)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) DeleteEvent(ctx context.Context, eventGUID uuid.UUID) error {
	sql := `
		DELETE FROM events
		WHERE guid = $1
	`
	_, err := s.conn.Exec(ctx, sql, eventGUID)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) FindEventsByInterval(ctx context.Context, startDateTime, endDateTime time.Time) ([]*storage.Event, error) {
	sql := `
		SELECT guid, title, start_at, end_at, description, user_guid, notify_before, notified
		FROM events
		WHERE start_at >= $1 AND end_at <= $2
	`

	var events []*storage.Event

	err := pgxscan.Select(ctx, s.conn, &events, sql, startDateTime, endDateTime)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s Storage) FindEventByGUID(ctx context.Context, eventGUID uuid.UUID) (*storage.Event, error) {
	query := `
		SELECT guid, title, start_at, end_at, description, user_guid, notify_before, notified
		FROM events
		WHERE guid = $1
	`

	var event storage.Event

	err := pgxscan.Get(ctx, s.conn, &event, query, eventGUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &event, nil
}

func (s *Storage) FindEventsNeedsNotify(ctx context.Context) ([]*storage.Event, error) {
	sql := `
		SELECT guid, title, start_at, end_at, description, user_guid, notify_before, notified
		FROM events
		WHERE notified = false AND start_at - notify_before <= $1
	`

	var events []*storage.Event

	err := pgxscan.Select(ctx, s.conn, &events, sql, time.Now())
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) DeleteEventsOlderThan(ctx context.Context, datetime time.Time) (rowAffected int64, err error) {
	sql := `
		DELETE FROM events
		WHERE start_at <= $1
	`

	commandTag, err := s.conn.Exec(ctx, sql, datetime)
	if err != nil {
		return 0, err
	}

	return commandTag.RowsAffected(), nil
}
