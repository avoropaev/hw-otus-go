package service

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	internalapp "github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/server/pb"
	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/storage"
)

type eventService struct {
	pb.UnimplementedEventServiceServer
	app internalapp.Application
}

func NewEventService(app internalapp.Application) pb.EventServiceServer {
	return &eventService{
		app: app,
	}
}

func (e *eventService) CreateEvent(ctx context.Context, event *pb.Event) (*pb.CreateUpdateResponse, error) {
	guid, err := uuid.Parse(event.GetGuid())
	if err != nil {
		return nil, err
	}

	userGUID, err := uuid.Parse(event.GetUserGuid())
	if err != nil {
		return nil, err
	}

	notifyBefore := event.GetNotifyBefore().AsDuration()

	err = e.app.CreateEvent(ctx, storage.Event{
		GUID:         guid,
		Title:        event.GetTitle(),
		StartAt:      event.GetStartAt().AsTime(),
		EndAt:        event.GetEndAt().AsTime(),
		Description:  event.Description,
		UserGUID:     userGUID,
		NotifyBefore: &notifyBefore,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateUpdateResponse{
		Event: event,
	}, nil
}

func (e *eventService) UpdateEvent(ctx context.Context, event *pb.Event) (*pb.CreateUpdateResponse, error) {
	guid, err := uuid.Parse(event.GetGuid())
	if err != nil {
		return nil, err
	}

	userGUID, err := uuid.Parse(event.GetUserGuid())
	if err != nil {
		return nil, err
	}

	notifyBefore := event.GetNotifyBefore().AsDuration()

	err = e.app.UpdateEvent(ctx, guid, storage.Event{
		GUID:         guid,
		Title:        event.GetTitle(),
		StartAt:      event.GetStartAt().AsTime(),
		EndAt:        event.GetEndAt().AsTime(),
		Description:  event.Description,
		UserGUID:     userGUID,
		NotifyBefore: &notifyBefore,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateUpdateResponse{
		Event: event,
	}, nil
}

func (e *eventService) DeleteEvent(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	guid, err := uuid.Parse(req.EventGuid)
	if err != nil {
		return nil, err
	}

	err = e.app.DeleteEvent(ctx, guid)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteResponse{}, nil
}

func (e *eventService) GetEventsForDay(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	events, err := e.app.GetEventForDay(ctx, req.GetStartAt().AsTime())
	if err != nil {
		return nil, err
	}

	return generateGetEventsResponse(events), nil
}

func (e *eventService) GetEventsForWeek(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	events, err := e.app.GetEventForWeek(ctx, req.GetStartAt().AsTime())
	if err != nil {
		return nil, err
	}

	return generateGetEventsResponse(events), nil
}

func (e *eventService) GetEventsForMonth(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	events, err := e.app.GetEventForMonth(ctx, req.GetStartAt().AsTime())
	if err != nil {
		return nil, err
	}

	return generateGetEventsResponse(events), nil
}

func generateGetEventsResponse(events []*storage.Event) *pb.GetEventsResponse {
	responseEvents := make([]*pb.Event, 0, len(events))

	for _, event := range events {
		responseEvents = append(responseEvents, &pb.Event{
			Guid:         event.GUID.String(),
			Title:        event.Title,
			StartAt:      &timestamppb.Timestamp{Seconds: event.StartAt.Unix()},
			EndAt:        &timestamppb.Timestamp{Seconds: event.EndAt.Unix()},
			Description:  event.Description,
			UserGuid:     event.UserGUID.String(),
			NotifyBefore: &durationpb.Duration{Seconds: int64(event.NotifyBefore.Seconds())},
		})
	}

	return &pb.GetEventsResponse{
		Events: responseEvents,
	}
}
