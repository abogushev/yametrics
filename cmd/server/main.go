// Сервер для получения и хранения метрик с агента
package main

import (
	"context"
	"log"
	_ "net/http/pprof"
	"os/signal"
	"syscall"
	"yametrics/internal/server"
	"yametrics/internal/server/config"
	"yametrics/internal/server/storage"

	"go.uber.org/zap"
)
//запуск сервера: инициализация хранилища и сервера
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("error on create logger: %v", err)
	}
	logger := l.Sugar()
	defer logger.Sync()

	cfgProvider := config.NewConfigProvider()

	var metricstorage storage.MetricsStorage

	if len(cfgProvider.StorageCfg.DBURL) != 0 {
		metricstorage, err = storage.NewDBMetricStorage(cfgProvider.StorageCfg.DBURL, ctx, logger)
	} else {
		metricstorage, err = storage.NewFileMetricsStorage(cfgProvider.StorageCfg, logger, ctx)
	}
	logger.Info("storage started successful")
	if err != nil {
		logger.Fatalf("error on create metric storage %v", err)
	}

	server.Run(logger, cfgProvider.ServerCfg, metricstorage, ctx)
	metricstorage.Close()
}
