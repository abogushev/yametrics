package main

import (
	"log"
	"net/http"
	"yametrics/cmd/server/handlers"
	"yametrics/cmd/server/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	metricsStorage := storage.NewGuageStorage()
	countersStorage := storage.NewCounterStorage()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/update", func(r chi.Router) {
		r.Post("/gauge/{name}/{value}", handlers.PostGuage(metricsStorage))
		r.Post("/counter/{name}/{value}", handlers.PostCounter(countersStorage))
		r.Post("/gauge/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })
		r.Post("/counter/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })
		r.Post("/*", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotImplemented) })
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", handlers.GetMetric(metricsStorage, countersStorage))
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.GetAllAsHTML(metricsStorage, countersStorage))
	})

	log.Println("server started successfull")
	log.Fatal(http.ListenAndServe(":8080", r))
}
