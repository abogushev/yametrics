package main

import (
	"context"
	"os/signal"
	"syscall"
	"yametrics/internal/agent"

	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	l, _ := zap.NewProduction()
	defer l.Sync() // flushes buffer, if any
	defer cancel()

	agent.Run(ctx, l.Sugar())

	<-ctx.Done()
}
