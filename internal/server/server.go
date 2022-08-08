package server

import (
	"context"
	"net/http"
	"yametrics/internal/server/config"
	"yametrics/internal/server/handlers"
	"yametrics/internal/server/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

func Run(
	logger *zap.SugaredLogger,
	cfg *config.ServerConfig,
	storage storage.MetricsStorage,
	dbstorage storage.DbMetricStorage,
	ctx context.Context) {
	handler := handlers.NewHandler(logger, storage, dbstorage, cfg.SignKey)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	r.Route("/update", func(r chi.Router) {
		r.Post("/{type:gauge|counter}/{name}/{value}", handler.UpdateV1)
		r.Post("/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotImplemented) })
		r.Post("/", handler.UpdateV2)
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", handler.GetV1)
		r.Post("/", handler.GetV2)
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/ping", handler.PingDb)
		r.Get("/", handler.GetAllAsHTML)
	})

	server := &http.Server{Addr: cfg.Address, Handler: r}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("server start error: %s\n", err)
		}
	}()
	logger.Info("server started successfuly")

	<-ctx.Done()
	logger.Info("get stop signal, start shutdown server")
	if err := server.Shutdown(ctx); err != nil && err != context.Canceled {
		logger.Fatalf("Server Shutdown Failed:%+v", err)
	} else {
		logger.Info("server stopped successfully")
	}
}
