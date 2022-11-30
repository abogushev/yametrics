// Агент сбора метрик
package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"sync"
	"syscall"
	"yametrics/internal/crypto"

	"go.uber.org/zap"

	"yametrics/internal/agent/config"
	"yametrics/internal/agent/managers"
	"yametrics/internal/metainfo"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// при старте запускаются два фоновых демона для сбора и отправки метрик на сервер с предопределенными конфигами.
func main() {
	metainfo.PrintBuildInfo(buildVersion, buildDate, buildCommit)

	l, _ := zap.NewProduction()
	logger := l.Sugar()
	defer logger.Sync()

	configProvider := config.NewConfigProvider()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	wg := &sync.WaitGroup{}

	publicKey, err := crypto.ReadPublicKey(configProvider.AgentCfg.CryptoKeyPath)
	if err != nil {
		logger.Errorf("init failed %v", err)
	}

	m := managers.NewMetricManager(logger, configProvider.AgentCfg)
	t := managers.NewTransportManager(logger, configProvider.AgentCfg, publicKey)

	m.RunAsync(ctx, wg)
	t.RunAsync(m.NotifyCh, ctx, wg)

	go func() {
		if err := http.ListenAndServe(":8100", nil); err != nil {
			logger.Fatalf("can't start metric server, %v", err)
		}
	}()

	wg.Wait()
}
