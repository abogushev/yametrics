package main

import (
	"yametrics/internal/server"
	"yametrics/internal/server/models"

	"github.com/caarlos0/env"
	"go.uber.org/zap"
)

func main() {
	l, _ := zap.NewProduction()
	logger := l.Sugar()
	defer logger.Sync()

	var cfg models.ServerConfig
	err := env.Parse(&cfg)
	if err != nil {
		logger.Fatal(err)
	}

	server.Run(logger, cfg)
}
