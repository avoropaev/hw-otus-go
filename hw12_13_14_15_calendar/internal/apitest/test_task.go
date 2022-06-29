package apitest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/gen/openapicli"
	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/storage"
)

func (s *APISuite) TestCreateEvent() {
	startAt := time.Date(2022, 3, 3, 6, 6, 6, 0, time.UTC)

	event := openapicli.EventEvent{
		Guid:         uuid.New().String(),
		Title:        "test-event",
		StartAt:      startAt,
		EndAt:        startAt.Add(time.Hour * 2),
		Description:  "test-event description",
		UserGuid:     uuid.New().String(),
		NotifyBefore: strconv.Itoa(int((time.Hour).Seconds())) + "s",
	}

	s.NewEvent(event)
	defer s.DeleteEvent(event.Guid)

	events := s.EventServiceGetEventsForDay(startAt)
	s.Require().NotNil(events)
	s.Require().Len(events, 1)
	s.Require().Equal(event, events[0])
}

func (s *APISuite) TestUpdateEvent() {
	event := openapicli.EventEvent{
		Guid:         uuid.New().String(),
		Title:        "test-event",
		StartAt:      time.Now().UTC(),
		EndAt:        time.Now().Add(time.Hour * 2).UTC(),
		Description:  "test-event description 1",
		UserGuid:     uuid.New().String(),
		NotifyBefore: strconv.Itoa(int((time.Hour).Seconds())) + "s",
	}

	s.NewEvent(event)

	event.Title = "test-event updated"

	respEvent, resp, err := s.cli.EventServiceApi.EventServiceUpdateEvent(s.ctx, event)
	s.Require().NoError(err, s.errBody(err))
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer func() {
		err := resp.Body.Close()
		s.Require().NoError(err, s.errBody(err))
	}()

	s.Require().Equal(event, respEvent.Event)

	s.DeleteEvent(event.Guid)
}

func (s *APISuite) TestDeleteEvent() {
	startAt := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	event := openapicli.EventEvent{
		Guid:         uuid.New().String(),
		Title:        "test-event",
		StartAt:      startAt,
		EndAt:        startAt.Add(time.Hour * 2).UTC(),
		Description:  "test-event description 1",
		UserGuid:     uuid.New().String(),
		NotifyBefore: strconv.Itoa(int((time.Hour).Seconds())) + "s",
	}

	s.NewEvent(event)

	events := s.EventServiceGetEventsForDay(startAt)
	s.Require().NotNil(events)
	s.Require().Len(events, 1)

	s.DeleteEvent(event.Guid)

	events = s.EventServiceGetEventsForDay(startAt)
	s.Require().NotNil(events)
	s.Require().Empty(events)
}

func (s *APISuite) TestGetEvents() {
	baseStartAt := time.Unix(5000, 0).UTC()
	baseEndAt := time.Unix(10000, 0).UTC()

	event := openapicli.EventEvent{
		Title:        "test",
		UserGuid:     uuid.New().String(),
		NotifyBefore: strconv.Itoa(int((time.Minute).Seconds())) + "s",
	}

	// добавляем 30 эвентов с интервалом 1 день
	for i := 0; i < 30; i++ {
		event.Guid = uuid.New().String()
		event.StartAt = baseStartAt.Add(time.Hour * 24 * time.Duration(i))
		event.EndAt = baseEndAt.Add(time.Hour * 24 * time.Duration(i))

		s.NewEvent(event)
	}

	events := s.EventServiceGetEventsForDay(time.Unix(4999, 0).UTC())
	s.Require().Len(events, 1)

	events = s.EventServiceGetEventsForWeek(time.Unix(4999, 0).UTC())
	s.Require().Len(events, 7)

	events = s.EventServiceGetEventsForMonth(time.Unix(4999, 0).UTC())
	s.Require().Len(events, 30)

	for _, event := range events {
		s.DeleteEvent(event.Guid)
	}
}

func (s *APISuite) TestUpdateEventErrEventNotFound() {
	event := openapicli.EventEvent{
		Guid:         uuid.New().String(),
		Title:        "test-event",
		StartAt:      time.Now().UTC(),
		EndAt:        time.Now().Add(time.Hour * 2).UTC(),
		Description:  "test-event description 1",
		UserGuid:     uuid.New().String(),
		NotifyBefore: strconv.Itoa(int((time.Hour).Seconds())) + "s",
	}

	_, resp, err := s.cli.EventServiceApi.EventServiceUpdateEvent(s.ctx, event)
	s.Require().Error(err)
	s.Require().Equal(http.StatusInternalServerError, resp.StatusCode)
	defer func() {
		err := resp.Body.Close()
		s.Require().NoError(err, s.errBody(err))
	}()

	s.Require().Equal(storage.ErrEventNotFound.Error(), s.errMessage(err))
}

func (s *APISuite) TestCreateEventErrEventAlreadyExists() {
	startAt := time.Date(2022, 3, 6, 6, 6, 6, 0, time.UTC)

	event := openapicli.EventEvent{
		Guid:         uuid.New().String(),
		Title:        "test-event",
		StartAt:      startAt,
		EndAt:        startAt.Add(time.Hour * 2),
		Description:  "test-event description",
		UserGuid:     uuid.New().String(),
		NotifyBefore: strconv.Itoa(int((time.Hour).Seconds())) + "s",
	}

	s.NewEvent(event)

	_, resp, err := s.cli.EventServiceApi.EventServiceCreateEvent(s.ctx, event)
	s.Require().Error(err)
	s.Require().Equal(http.StatusInternalServerError, resp.StatusCode)
	defer func() {
		err := resp.Body.Close()
		s.Require().NoError(err, s.errBody(err))
	}()

	s.Require().Equal(storage.ErrEventAlreadyExists.Error(), s.errMessage(err))
}
