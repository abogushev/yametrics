package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"yametrics/internal/server"
	"yametrics/internal/server/models"
	"yametrics/internal/server/storage"

	"github.com/caarlos0/env"
	"go.uber.org/zap"
)

func main() {
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal("error on create logger", err)
	}
	logger := l.Sugar()
	defer logger.Sync()

	var serverCfg models.ServerConfig
	if err = env.Parse(&serverCfg); err != nil {
		logger.Fatal("error on read cfg of server", err)
	}

	var storageCfg models.MetricsStorageConfig
	if err = env.Parse(&storageCfg); err != nil {
		logger.Fatal("error on read cfg of metric storage", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	var metricsStorage storage.MetricsStorage
	if metricsStorage, err = storage.NewMetricsStorageImpl(&storageCfg, logger, ctx); err != nil {
		logger.Fatal("error on create metric storage", err)
	}

	server.Run(logger, serverCfg, metricsStorage, ctx)
}
