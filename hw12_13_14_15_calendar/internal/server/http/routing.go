// Package internalhttp implement work with http protocol
package internalhttp

import (
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// MakeRouter creates handler for http with all routes
func MakeRouter(app Application) http.Handler {
	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc("/liveness", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte("alive"))
		if err != nil {
			log.Error().Err(err)
		}
	})

	r.HandleFunc("/test-db", func(writer http.ResponseWriter, request *http.Request) {
		err := app.TestDB(request.Context())
		if err != nil {
			log.Error().Err(err).Send()
		}
	})

	r.Use(loggingMiddleware())

	registerPprof(r)

	return r
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
