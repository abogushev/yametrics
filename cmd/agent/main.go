package main

import (
	"context"
	"os/signal"
	"syscall"
	"yametrics/internal/agent"
	"yametrics/internal/agent/config"

	"go.uber.org/zap"
)

func main() {
	l, _ := zap.NewProduction()
	logger := l.Sugar()
	defer logger.Sync()

	configProvider := config.NewConfigProvider()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	agent.NewAgent(logger, configProvider.AgentCfg).RunSync(ctx)
}
