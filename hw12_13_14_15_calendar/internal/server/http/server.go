package internalhttp

import (
	"context"
	golog "log"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/server"
	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/pkg/zerologwriter"
)

type serv struct {
	host         string
	port         int
	grpcEndpoint string
	app          app.Application

	server *http.Server
}

var _ server.IServer = (*serv)(nil)

func NewServer(host string, port int, grpcEndpoint string, app app.Application) server.IServer {
	return &serv{host, port, grpcEndpoint, app, nil}
}

func (s *serv) Start(ctx context.Context) error {
	handler, err := MakeRouter(ctx, s.grpcEndpoint, s.app)
	if err != nil {
		return err
	}

	s.server = &http.Server{
		Addr:         s.host + ":" + strconv.Itoa(s.port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		ErrorLog:     golog.New(zerologwriter.ZerologWriter{Zerolog: log.Logger}, "", golog.LstdFlags),
	}

	return s.server.ListenAndServe()
}

func (s *serv) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	return s.server.Shutdown(ctx)
}
