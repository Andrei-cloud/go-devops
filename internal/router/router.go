package router

import (
	"github.com/andrei-cloud/go-devops/internal/handlers"
	mw "github.com/andrei-cloud/go-devops/internal/middlewares"
	"github.com/andrei-cloud/go-devops/internal/repo"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

func SetupRouter(repo repo.Repository, key []byte) *chi.Mux {
	log.Debug().Msg("Setting up the router")
	r := chi.NewRouter()
	//r.Use(middleware.Logger)
	r.Use(mw.GzipMW, mw.KeyInject(key))
	r.Get("/", handlers.Default())
	r.Get("/value/{m_type}/{m_name}", handlers.GetMetrics(repo))
	r.Get("/ping", handlers.Ping(repo))

	r.Post("/update/{m_type}/{m_name}/{value}", handlers.Update(repo))
	r.Post("/update/", handlers.UpdatePost(repo))
	r.Post("/updates/", handlers.UpdateBulkPost(repo))
	r.Post("/value/", handlers.GetMetricsPost(repo))

	return r
}
