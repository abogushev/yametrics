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
		log.Fatal("error on create logger", err)
	}
	logger := l.Sugar()
	defer logger.Sync()

	cfgProvider := config.NewConfigProvider()

	var metricstorage storage.MetricsStorage

	if dbUrl := cfgProvider.StorageCfg.DbUrl; dbUrl != "" {
		metricstorage, err = storage.NewDbMetricStorage(dbUrl, ctx)
	} else {
		metricstorage, err = storage.NewFileMetricsStorage(cfgProvider.StorageCfg, logger, ctx)
	}

	if err != nil {
		logger.Fatal("error on create metric storage", err)
	}

	server.Run(logger, cfgProvider.ServerCfg, metricstorage, ctx)
	metricstorage.Close()
}
