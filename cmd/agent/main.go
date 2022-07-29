package main

import (
	"context"
	"os/signal"
	"syscall"
	"yametrics/internal/agent"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()
	agent.NewAgent().RunSync(ctx)
}
