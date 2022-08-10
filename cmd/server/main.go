package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"yametrics/internal/server"
	"yametrics/internal/server/config"
	"yametrics/internal/server/storage"

	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("error on create logger: %w", err)
	}
	logger := l.Sugar()
	defer logger.Sync()

	cfgProvider := config.NewConfigProvider()

	var metricstorage storage.MetricsStorage

	if dbURL := cfgProvider.StorageCfg.DBURL; dbURL != "" {
		metricstorage, err = storage.NewDBMetricStorage(dbURL, ctx)
	} else {
		metricstorage, err = storage.NewFileMetricsStorage(cfgProvider.StorageCfg, logger, ctx)
	}
	logger.Info("storage started successful")
	if err != nil {
		logger.Fatalf("error on create metric storage %w", err)
	}

	server.Run(logger, cfgProvider.ServerCfg, metricstorage, ctx)
	metricstorage.Close()
}
