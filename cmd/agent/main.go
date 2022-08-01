package main

import (
	"context"
	"os/signal"
	"syscall"
	"yametrics/internal/agent"
	"yametrics/internal/agent/models"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

func main() {
	l, _ := zap.NewProduction()
	logger := l.Sugar()
	defer logger.Sync()

	var cfg models.AgentConfig
	err := env.Parse(&cfg)
	if err != nil {
		logger.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	agent.NewAgent(logger, cfg).RunSync(ctx)
}
