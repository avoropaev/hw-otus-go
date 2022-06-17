package memorystorage_test

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/storage/memory"
)

type fixture struct {
	Name    string
	Start   time.Time
	End     time.Time
	MustGot bool
}

func TestStorage_Create_Update_Delete(t *testing.T) {
	eventStorage := memorystorage.New()
	ctx := context.Background()

	// CREATE
	description := "description"
	notifyBefore := time.Hour

	event := storage.Event{
		GUID:         uuid.New(),
		Title:        "Test",
		StartAt:      time.Now(),
		EndAt:        time.Now().Add(time.Hour * 1),
		Description:  &description,
		UserGUID:     uuid.New(),
		NotifyBefore: &notifyBefore,
	}

	err := eventStorage.CreateEvent(ctx, event)
	require.NoError(t, err)

	resultEvent, err := eventStorage.FindEventByGUID(ctx, event.GUID)
	require.NoError(t, err)
	require.NotNil(t, resultEvent)
	require.Equal(t, &event, resultEvent)

	err = eventStorage.CreateEvent(ctx, event)
	require.ErrorIs(t, err, storage.ErrEventAlreadyExists)

	// UPDATE
	newEvent := storage.Event{
		GUID:     uuid.New(),
		StartAt:  time.Now().Add(time.Hour * 10),
		EndAt:    time.Now().Add(time.Hour * 15),
		UserGUID: uuid.New(),
	}

	err = eventStorage.UpdateEvent(ctx, event.GUID, newEvent)
	require.NoError(t, err)

	newEvent.GUID = event.GUID

	resultEvent, err = eventStorage.FindEventByGUID(ctx, event.GUID)
	require.NoError(t, err)
	require.NotNil(t, resultEvent)
	require.Equal(t, &newEvent, resultEvent)

	err = eventStorage.UpdateEvent(ctx, uuid.New(), newEvent)
	require.ErrorIs(t, err, storage.ErrEventNotFound)

	// DELETE
	err = eventStorage.DeleteEvent(ctx, event.GUID)
	require.NoError(t, err)

	resultEvent, err = eventStorage.FindEventByGUID(ctx, event.GUID)
	require.NoError(t, err)
	require.Nil(t, resultEvent)

	err = eventStorage.DeleteEvent(ctx, event.GUID)
	require.NoError(t, err)
}

func TestStorage_FindEventsByInterval(t *testing.T) {
	iStart := time.Now()
	iEnd := iStart.AddDate(0, 0, 1)
	eventStorage := memorystorage.New()

	ctx := context.Background()
	expected := make([]*storage.Event, 0)

	fixtures := []fixture{
		{Name: "whole before interval", Start: iStart.Add(time.Hour * -5), End: iStart.Add(time.Hour * -3), MustGot: false},
		{Name: "start before iStart and end in iStart", Start: iStart.Add(time.Hour * -2), End: iStart, MustGot: true},
		{Name: "start before iStart and end in interval", Start: iStart.Add(time.Hour * -2), End: iStart.Add(time.Hour * 2), MustGot: true},
		{Name: "start in iStart and end in iStart", Start: iStart, End: iStart, MustGot: true},
		{Name: "start in iStart and end in interval", Start: iStart, End: iStart.Add(time.Hour * 2), MustGot: true},
		{Name: "in interval", Start: iStart.Add(time.Hour * 2), End: iStart.Add(time.Hour * 4), MustGot: true},
		{Name: "start in interval and end in iEnd", Start: iStart.Add(time.Hour * 2), End: iEnd, MustGot: true},
		{Name: "start in iEnd and end in iEnd", Start: iEnd, End: iEnd, MustGot: true},
		{Name: "start in interval and end after iEnd", Start: iStart.Add(time.Hour * 2), End: iEnd.Add(time.Hour * 3), MustGot: true},
		{Name: "start in iEnd and end after iEnd", Start: iEnd, End: iEnd.Add(time.Hour * 3), MustGot: true},
		{Name: "whole after interval", Start: iEnd.Add(time.Hour * 3), End: iEnd.Add(time.Hour * 6), MustGot: false},
		{Name: "start before iStart and end after iEnd", Start: iStart.Add(time.Hour * -3), End: iEnd.Add(time.Hour * 3), MustGot: true},
	}

	for _, fixture := range fixtures {
		event := storage.Event{
			GUID:    uuid.New(),
			Title:   fixture.Name,
			StartAt: fixture.Start,
			EndAt:   fixture.End,
		}

		err := eventStorage.CreateEvent(ctx, event)
		require.NoError(t, err)

		if fixture.MustGot {
			expected = append(expected, &event)
		}
	}

	result, err := eventStorage.FindEventsByInterval(ctx, iStart, iEnd)
	require.NoError(t, err)
	equalSlices(t, expected, result)
}

func TestStorage_Concurrency(t *testing.T) {
	eventStorage := memorystorage.New()
	ctx := context.Background()

	// CREATE
	description := "description"
	notifyBefore := time.Hour

	event := storage.Event{
		GUID:         uuid.New(),
		Title:        "Test",
		StartAt:      time.Now(),
		EndAt:        time.Now().Add(time.Hour * 1),
		Description:  &description,
		UserGUID:     uuid.New(),
		NotifyBefore: &notifyBefore,
	}

	wg := sync.WaitGroup{}

	events := make([]*storage.Event, 0, 10)
	for i := 0; i < 10; i++ {
		wg.Add(1)

		event := event
		event.GUID = uuid.New()
		event.Title = "Test " + strconv.Itoa(i)

		events = append(events, &event)

		go func() {
			defer wg.Done()

			err := eventStorage.CreateEvent(ctx, event)
			require.NoError(t, err)
		}()
	}

	wg.Wait()

	result, err := eventStorage.FindEventsByInterval(ctx, time.Now(), time.Now().Add(time.Hour*24))
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Len(t, result, len(events))
	equalSlices(t, events, result)

	// UPDATE
	newEvent := storage.Event{
		GUID:     uuid.New(),
		StartAt:  time.Now().Add(time.Hour * 10),
		EndAt:    time.Now().Add(time.Hour * 15),
		UserGUID: uuid.New(),
	}

	for i, event := range events {
		wg.Add(1)

		event := event
		events[i] = &storage.Event{
			GUID:     event.GUID,
			StartAt:  newEvent.StartAt,
			EndAt:    newEvent.EndAt,
			UserGUID: newEvent.UserGUID,
		}

		go func() {
			defer wg.Done()

			err := eventStorage.UpdateEvent(ctx, event.GUID, newEvent)
			require.NoError(t, err)
		}()
	}

	wg.Wait()

	result, err = eventStorage.FindEventsByInterval(ctx, time.Now(), time.Now().Add(time.Hour*24))
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Len(t, result, len(events))
	equalSlices(t, events, result)

	// DELETE
	for i := 0; i < 5; i++ {
		for _, event := range events {
			wg.Add(1)

			event := event

			go func() {
				defer wg.Done()

				err := eventStorage.DeleteEvent(ctx, event.GUID)
				require.NoError(t, err)
			}()
		}
	}

	wg.Wait()

	result, err = eventStorage.FindEventsByInterval(ctx, time.Now(), time.Now().Add(time.Hour*24))
	require.NoError(t, err)
	require.Empty(t, result)
}

func equalSlices(t *testing.T, expected, result []*storage.Event) {
	t.Helper()

	require.Equal(t, len(expected), len(result))

	for _, expectedEvent := range expected {
		var equalEvent *storage.Event

		for _, resultEvent := range result {
			if expectedEvent.GUID == resultEvent.GUID {
				equalEvent = resultEvent
				break
			}
		}

		require.Equal(t, expectedEvent, equalEvent)
	}
}
