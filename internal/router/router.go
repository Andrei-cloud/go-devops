// Pakage provides functions to setup router for http server.
package router

import (
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"

	"github.com/andrei-cloud/go-devops/internal/encrypt"
	"github.com/andrei-cloud/go-devops/internal/handlers"
	mw "github.com/andrei-cloud/go-devops/internal/middlewares"
	"github.com/andrei-cloud/go-devops/internal/repo"
)

// SetupRouter -  Function setup chi router for handdlers and required middlewares
//     repo - take entity implementing Repository interface
//     key - slice of bytes of key for hash validation.
func SetupRouter(repo repo.Repository, key []byte, e encrypt.Encrypter) *chi.Mux {
	log.Debug().Msg("Setting up the router")
	r := chi.NewRouter()
	r.Use(mw.CryptoMW(e), mw.GzipMW, mw.KeyInject(key))
	r.Get("/", handlers.Default())
	r.Get("/value/{m_type}/{m_name}", handlers.GetMetrics(repo))
	r.Get("/ping", handlers.Ping(repo))

	r.Post("/update/{m_type}/{m_name}/{value}", handlers.Update(repo))
	r.Post("/update/", handlers.UpdatePost(repo))
	r.Post("/updates/", handlers.UpdateBulkPost(repo))
	r.Post("/value/", handlers.GetMetricsPost(repo))

	return r
}

// WithPPROF - Function to setup router for PPROF handlers
//    r tange chu router to enrach with pprof handlers.
func WithPPROF(r *chi.Mux) *chi.Mux {
	r.Handle("/debug/pprof", http.HandlerFunc(pprof.Index))
	r.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	r.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	r.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))
	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))

	return r
}
