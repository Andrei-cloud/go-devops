package router

import (
	"github.com/andrei-cloud/go-devops/internal/handlers"
	"github.com/andrei-cloud/go-devops/internal/storage/inmem"
	"github.com/go-chi/chi"
)

func SetupRouter() *chi.Mux {

	repo := inmem.New()

	r := chi.NewRouter()
	//r.Use(middleware.Logger)
	//r.Use(middleware.Recoverer)
	r.Post("/update/counter/{m_name}/{value}", handlers.Counters(repo))
	r.Post("/update/gauge/{m_name}/{value}", handlers.Gauges(repo))

	return r
}
