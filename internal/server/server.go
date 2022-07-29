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
	handler := handlers.NewHandler(l.Sugar(), storage.NewMetricsStorageImpl())

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/update", func(r chi.Router) {
		r.Post("/{type:gauge|counter}/{name}/{value}", handler.UpdateV1)
		r.Post("/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotImplemented) })
		// r.Post("/*", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotImplemented) })
		r.Post("/", handler.UpdateV2)
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", handler.GetV1)
		r.Post("/", handler.GetV2)
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", handler.GetAllAsHTML)
	})

	log.Println("server started successfull")
	log.Fatal(http.ListenAndServe(":8080", r))
}
