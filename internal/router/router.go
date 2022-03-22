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
	r.Get("/value/{m_type}/{m_name}", handlers.GetMetrics(repo))
	r.Post("/update/{m_type}/{m_name}/{value}", handlers.Update(repo))
	// r.Post("/update/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.Error(w, "invalid metric type", http.StatusNotImplemented)
	// })

	return r
}
