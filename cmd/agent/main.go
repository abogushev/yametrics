package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"
	"yametrics/internal/agent/config"
	"yametrics/internal/agent/managers"

	"go.uber.org/zap"
)

func main() {
	l, _ := zap.NewProduction()
	logger := l.Sugar()
	defer logger.Sync()

	configProvider := config.NewConfigProvider()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	wg := &sync.WaitGroup{}

	m := managers.NewMetricManager(logger, configProvider.AgentCfg)
	t := managers.NewTransportManager(logger, configProvider.AgentCfg)

	m.RunAsync(ctx, wg)
	t.RunAsync(m.NotifyCh, ctx, wg)

	wg.Wait()
}
