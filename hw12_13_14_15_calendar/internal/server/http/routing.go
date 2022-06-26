// Package internalhttp implement work with http protocol
package internalhttp

import (
	"context"
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/server/pb"
)

// MakeRouter creates handler for http with all routes.
func MakeRouter(ctx context.Context, grpcEndpoint string, _ app.Application) (http.Handler, error) {
	r := mux.NewRouter()
	r.StrictSlash(true)

	registerPprof(r)

	r.HandleFunc("/liveness", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte("alive"))
		if err != nil {
			log.Error().Err(err)
		}
	})

	r.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		http.ServeFile(w, r, "gen/openapiv2/EventService.swagger.json")
	})

	apiHandler := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterEventServiceHandlerFromEndpoint(ctx, apiHandler, grpcEndpoint, opts)
	if err != nil {
		return nil, err
	}

	r.PathPrefix("/").Handler(apiHandler)

	r.Use(loggingMiddleware())

	return r, nil
}

func registerPprof(r *mux.Router) {
	s := r.PathPrefix("/pprof").Subrouter()
	s.HandleFunc("/", pprof.Index)
	s.HandleFunc("/cmdline", pprof.Cmdline)
	s.HandleFunc("/profile", pprof.Profile)
	s.HandleFunc("/symbol", pprof.Symbol)
	s.HandleFunc("/trace", pprof.Trace)
	s.Handle("/allocs", pprof.Handler("allocs"))
	s.Handle("/block", pprof.Handler("block"))
	s.Handle("/goroutine", pprof.Handler("goroutine"))
	s.Handle("/heap", pprof.Handler("heap"))
	s.Handle("/mutex", pprof.Handler("mutex"))
	s.Handle("/threadcreate", pprof.Handler("threadcreate"))
}
