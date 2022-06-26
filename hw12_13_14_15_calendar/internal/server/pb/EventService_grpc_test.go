package pb_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/server/grpc/service"
	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/server/pb"
	memorystorage "github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/storage/memory"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

type EventServiceGRPCSuite struct {
	suite.Suite
	Server *grpc.Server
	Conn   *grpc.ClientConn
	Client pb.EventServiceClient
}

func (suite *EventServiceGRPCSuite) SetupSuite() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterEventServiceServer(s, service.NewEventService(app.New(memorystorage.New())))
	go func() {
		if err := s.Serve(lis); err != nil {
			require.NoError(suite.T(), err)
		}
	}()

	ctx := context.Background()

	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(suite.T(), err)

	suite.Server = s
	suite.Conn = conn
	suite.Client = pb.NewEventServiceClient(conn)
}

func (suite *EventServiceGRPCSuite) TearDownSuite() {
	err := suite.Conn.Close()
	require.NoError(suite.T(), err)

	suite.Server.GracefulStop()
}

func (suite *EventServiceGRPCSuite) SetupTest() {
	events, err := suite.Client.GetEventsForMonth(context.Background(), &pb.GetEventsRequest{StartAt: &timestamppb.Timestamp{Seconds: 499}})
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), events)

	for _, event := range events.Events {
		result, err := suite.Client.DeleteEvent(context.Background(), &pb.DeleteRequest{EventGuid: event.Guid})
		require.NoError(suite.T(), err)
		require.Nil(suite.T(), result.Error)
	}
}

func TestEventServiceGRPCSuite(t *testing.T) {
	suite.Run(t, new(EventServiceGRPCSuite))
}

func (suite *EventServiceGRPCSuite) TestCreateEvent() {
	event := &pb.Event{
		Guid:         uuid.New().String(),
		Title:        "title",
		StartAt:      &timestamppb.Timestamp{Seconds: 500},
		EndAt:        &timestamppb.Timestamp{Seconds: 1000},
		UserGuid:     uuid.New().String(),
		NotifyBefore: &durationpb.Duration{Seconds: 600},
	}

	resp, err := suite.Client.CreateEvent(context.Background(), event)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), event.String(), resp.GetEvent().String())

	events, err := suite.Client.GetEventsForDay(context.Background(), &pb.GetEventsRequest{StartAt: &timestamppb.Timestamp{Seconds: 400}})
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), events)
	require.Equal(suite.T(), 1, len(events.Events))
	require.Equal(suite.T(), event.String(), events.Events[0].String())
}

func (suite *EventServiceGRPCSuite) TestUpdateEvent() {
	event := &pb.Event{
		Guid:         uuid.New().String(),
		Title:        "title 2",
		StartAt:      &timestamppb.Timestamp{Seconds: 500},
		EndAt:        &timestamppb.Timestamp{Seconds: 1000},
		UserGuid:     uuid.New().String(),
		NotifyBefore: &durationpb.Duration{Seconds: 60},
	}

	resp, err := suite.Client.CreateEvent(context.Background(), event)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), event.String(), resp.GetEvent().String())

	event.Title = "title 2 updated"

	resp, err = suite.Client.UpdateEvent(context.Background(), event)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), event.String(), resp.GetEvent().String())

	events, err := suite.Client.GetEventsForDay(context.Background(), &pb.GetEventsRequest{StartAt: &timestamppb.Timestamp{Seconds: 400}})
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), events)
	require.Equal(suite.T(), 1, len(events.Events))
	require.Equal(suite.T(), event.String(), events.Events[0].String())
}

func (suite *EventServiceGRPCSuite) TestDeleteEvent() {
	event := &pb.Event{
		Guid:         uuid.New().String(),
		Title:        "title 3",
		StartAt:      &timestamppb.Timestamp{Seconds: 500},
		EndAt:        &timestamppb.Timestamp{Seconds: 1000},
		UserGuid:     uuid.New().String(),
		NotifyBefore: &durationpb.Duration{Seconds: 60},
	}

	resp, err := suite.Client.CreateEvent(context.Background(), event)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), event.String(), resp.GetEvent().String())

	result, err := suite.Client.DeleteEvent(context.Background(), &pb.DeleteRequest{EventGuid: event.Guid})
	require.NoError(suite.T(), err)
	require.Nil(suite.T(), result.Error)

	events, err := suite.Client.GetEventsForDay(context.Background(), &pb.GetEventsRequest{StartAt: &timestamppb.Timestamp{Seconds: 400}})
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), events)
	require.Equal(suite.T(), 0, len(events.Events))
}

func (suite *EventServiceGRPCSuite) TestGetEvents() {
	baseStartAt := time.Unix(500, 0)
	baseEndAt := time.Unix(1000, 0)

	event := &pb.Event{
		Title:        "title 4",
		UserGuid:     uuid.New().String(),
		NotifyBefore: &durationpb.Duration{Seconds: 60},
	}

	// добавляем 30 эвентов с интервалом 1 день
	for i := 0; i < 30; i++ {
		event.Guid = uuid.New().String()
		event.StartAt = &timestamppb.Timestamp{Seconds: baseStartAt.Add(time.Hour * 24 * time.Duration(i)).Unix()}
		event.EndAt = &timestamppb.Timestamp{Seconds: baseEndAt.Add(time.Hour * 24 * time.Duration(i)).Unix()}

		resp, err := suite.Client.CreateEvent(context.Background(), event)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), event.String(), resp.GetEvent().String())
	}

	events, err := suite.Client.GetEventsForDay(context.Background(), &pb.GetEventsRequest{StartAt: &timestamppb.Timestamp{Seconds: 499}})
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), events)
	require.Equal(suite.T(), 1, len(events.Events))

	events, err = suite.Client.GetEventsForWeek(context.Background(), &pb.GetEventsRequest{StartAt: &timestamppb.Timestamp{Seconds: 499}})
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), events)
	require.Equal(suite.T(), 7, len(events.Events))

	events, err = suite.Client.GetEventsForMonth(context.Background(), &pb.GetEventsRequest{StartAt: &timestamppb.Timestamp{Seconds: 499}})
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), events)
	require.Equal(suite.T(), 30, len(events.Events))
}
