package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"yametrics/internal/server"
	"yametrics/internal/server/config"
	"yametrics/internal/server/storage"

	"go.uber.org/zap"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	wg := new(sync.WaitGroup)

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal("error on create logger", err)
	}
	logger := l.Sugar()
	defer logger.Sync()

	cfgProvider := config.NewConfigProvider()

	logger.Infof("server conf: %v", cfgProvider.ServerCfg)
	logger.Infof("storage conf: %v", cfgProvider.StorageCfg)

	var metricsStorage storage.MetricsStorage
	if metricsStorage, err = storage.NewMetricsStorageImpl(cfgProvider.StorageCfg, logger, ctx, wg); err != nil {
		logger.Fatal("error on create metric storage", err)
	}

	dbstorage := storage.NewDbMetricStorage(cfgProvider.StorageCfg.DbUrl, ctx)

	server.Run(logger, cfgProvider.ServerCfg, metricsStorage, *dbstorage, ctx)
	wg.Wait()
}
