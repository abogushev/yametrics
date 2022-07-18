package server

import (
	"log"
	"net/http"
	"yametrics/internal/server/handlers"
	"yametrics/internal/server/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

func Run() {
	l, _ := zap.NewProduction()
	defer l.Sync() // flushes buffer, if any
	logger := l.Sugar()

	metricsStorage := storage.NewGuageStorage()
	countersStorage := storage.NewCounterStorage()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/update", func(r chi.Router) {
		r.Post("/gauge/{name}/{value}", handlers.PostGuage(metricsStorage, logger))
		r.Post("/counter/{name}/{value}", handlers.PostCounter(countersStorage, logger))
		r.Post("/gauge/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })
		r.Post("/counter/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })
		r.Post("/*", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotImplemented) })
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", handlers.GetMetric(metricsStorage, countersStorage, logger))
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.GetAllAsHTML(metricsStorage, countersStorage, logger))
	})

	log.Println("server started successfull")
	log.Fatal(http.ListenAndServe(":8080", r))
}
