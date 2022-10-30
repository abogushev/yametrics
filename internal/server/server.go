package server

import (
	"context"
	"errors"
	"net/http"
	_ "net/http/pprof"
	"yametrics/internal/server/config"
	"yametrics/internal/server/handlers"
	"yametrics/internal/server/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)
//Run - запуск сервера
func Run(
	logger *zap.SugaredLogger,
	cfg *config.ServerConfig,
	storage storage.MetricsStorage,
	ctx context.Context) {
	handler := handlers.NewHandler(logger, storage, cfg.SignKey)

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

	r.Route("/updates", func(r chi.Router) {
		r.Post("/", handler.UpdatesV2)
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", handler.GetV1)
		r.Post("/", handler.GetV2)
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/ping", handler.PingDB)
		r.Get("/", handler.GetAllAsHTML)
	})

	runProfileServer(logger)

	server := &http.Server{Addr: cfg.Address, Handler: r}

	go func() {
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("server start error: %v", err)
		}
	}()
	logger.Infof("server started successfuly, addr:%v", server.Addr)

	<-ctx.Done()
	logger.Info("get stop signal, start shutdown server")
	if err := server.Shutdown(ctx); err != nil && errors.Is(err, context.Canceled) {
		logger.Fatalf("Server Shutdown Failed:%v", err)
	} else {
		logger.Info("server stopped successfully")
	}
}

func runProfileServer(logger *zap.SugaredLogger) {
	go func() {
		if err := http.ListenAndServe(":8200", nil); err != nil {
			logger.Fatalf("can't start metric server, %v", err)
		}
	}()
}
