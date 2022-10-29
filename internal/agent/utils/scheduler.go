package utils

import (
	"context"
	"time"

	"go.uber.org/zap"
)

func Schedule(f func(), ctx context.Context, duration time.Duration, name string, logger *zap.SugaredLogger) {
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			logger.Infof("call task: %s", name)
			f()

		case <-ctx.Done():
			ticker.Stop()
			logger.Infof("cancel task: %s", name)
			return
		}
	}
}
