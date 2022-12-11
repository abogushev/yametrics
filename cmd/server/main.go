// Сервер для получения и хранения метрик с агента
package main

import (
	"context"
	"log"
	_ "net/http/pprof"
	"os/signal"
	"syscall"
	"yametrics/internal/crypto"
	"yametrics/internal/server"

	"go.uber.org/zap"

	"yametrics/internal/metainfo"
	"yametrics/internal/server/config"
	"yametrics/internal/server/grpc"
	"yametrics/internal/server/storage"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// запуск сервера: инициализация хранилища и сервера
func main() {
	metainfo.PrintBuildInfo(buildVersion, buildDate, buildCommit)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("error on create logger: %v", err)
	}
	logger := l.Sugar()
	defer logger.Sync()

	cfg := config.NewServerConfig()

	var metricstorage storage.MetricsStorage

	if len(cfg.DBURL) != 0 {
		metricstorage, err = storage.NewDBMetricStorage(cfg.DBURL, ctx, logger)
	} else {
		metricstorage, err = storage.NewFileMetricsStorage(cfg, logger, ctx)
	}
	logger.Info("storage started successful")
	if err != nil {
		logger.Fatalf("error on create metric storage %v", err)
	}

	privateKey, err := crypto.ReadPrivateKey(cfg.CryptoKeyPath)
	if err != nil {
		logger.Errorf("error on read private key, %v", err)
	}

	go grpc.RunMetricsServer(logger, ctx)
	server.Run(logger, cfg, metricstorage, ctx, privateKey)
	metricstorage.Close()
}
