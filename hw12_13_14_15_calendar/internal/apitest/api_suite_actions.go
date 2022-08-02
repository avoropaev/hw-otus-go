package apitest

import (
	"context"
	"net/http"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/gen/openapicli"
)

type APISuiteActions struct {
	suite.Suite
	cli    *openapicli.APIClient
	ctx    context.Context
	apiURL string
}

func (s *APISuiteActions) Init(apiURL string) {
	apiCfg := openapicli.NewConfiguration()
	s.cli = openapicli.NewAPIClient(apiCfg)
	s.cli.ChangeBasePath(apiURL)
	s.ctx = context.Background()
	s.apiURL = apiURL
}

func (s *APISuiteActions) Client() *openapicli.EventServiceApiService {
	return s.cli.EventServiceApi
}

func (s *APISuite) NewEvent(event openapicli.EventEvent) {
	respEvent, resp, err := s.cli.EventServiceApi.EventServiceCreateEvent(s.ctx, event)
	s.Require().NoError(err, s.errBody(err))
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer func() {
		err := resp.Body.Close()
		s.Require().NoError(err, s.errBody(err))
	}()

	s.Require().Equal(event, respEvent.Event)
}

func (s *APISuite) DeleteEvent(eventGUID string) {
	respEvent, resp, err := s.cli.EventServiceApi.EventServiceDeleteEvent(s.ctx, openapicli.EventDeleteRequest{EventGuid: eventGUID})
	s.Require().NoError(err, s.errBody(err))
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer func() {
		err := resp.Body.Close()
		s.Require().NoError(err, s.errBody(err))
	}()

	s.Require().Empty(respEvent.Error)
}

func (s *APISuite) EventServiceGetEventsForDay(startAt time.Time) []openapicli.EventEvent {
	events, resp, err := s.cli.EventServiceApi.EventServiceGetEventsForDay(s.ctx, openapicli.EventGetEventsRequest{StartAt: startAt})
	s.Require().NoError(err, s.errBody(err))
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer func() {
		err := resp.Body.Close()
		s.Require().NoError(err, s.errBody(err))
	}()

	return events.Events
}

func (s *APISuite) EventServiceGetEventsForWeek(startAt time.Time) []openapicli.EventEvent {
	events, resp, err := s.cli.EventServiceApi.EventServiceGetEventsForWeek(s.ctx, openapicli.EventGetEventsRequest{StartAt: startAt})
	s.Require().NoError(err, s.errBody(err))
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer func() {
		err := resp.Body.Close()
		s.Require().NoError(err, s.errBody(err))
	}()

	return events.Events
}

func (s *APISuite) EventServiceGetEventsForMonth(startAt time.Time) []openapicli.EventEvent {
	events, resp, err := s.cli.EventServiceApi.EventServiceGetEventsForMonth(s.ctx, openapicli.EventGetEventsRequest{StartAt: startAt})
	s.Require().NoError(err, s.errBody(err))
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer func() {
		err := resp.Body.Close()
		s.Require().NoError(err, s.errBody(err))
	}()

	return events.Events
}
