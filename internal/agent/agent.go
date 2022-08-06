package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"yametrics/internal/agent/config"
	"yametrics/internal/agent/models/storage"
	"yametrics/internal/protocol"

	"go.uber.org/zap"
)

type Agent struct {
	url     string
	client  http.Client
	logger  *zap.SugaredLogger
	metrics *storage.Metrics
	config  *config.AgentConfig
}

func NewAgent(l *zap.SugaredLogger, config *config.AgentConfig) *Agent {
	return &Agent{
		url:     "http://" + config.Address,
		client:  http.Client{},
		logger:  l,
		metrics: &storage.Metrics{MemStats: &runtime.MemStats{}, PollCount: 0, RandomValue: 0.0},
		config:  config,
	}
}

func (agent *Agent) RunSync(ctx context.Context) {
	go agent.updateMetricsWithInterval(ctx)
	go agent.sendMetricsWithInterval(ctx)
	agent.logger.Info("agent started")
	<-ctx.Done()
	agent.logger.Info("agent stoped")
}

func (agent *Agent) schedule(f func(), ctx context.Context, duration time.Duration, name string) {
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			agent.logger.Infof("call task: %s", name)
			f()

		case <-ctx.Done():
			ticker.Stop()
			agent.logger.Infof("cancel task: %s", name)
			return
		}
	}
}

func (agent *Agent) updateMetricsWithInterval(ctx context.Context) {
	agent.schedule(
		func() {
			runtime.ReadMemStats(agent.metrics.MemStats)
			agent.metrics.PollCount++
			agent.metrics.RandomValue = rand.Float64()
		},
		ctx,
		agent.config.PollInterval,
		"collecting metrics")
}

func (agent *Agent) sendMetricsWithInterval(ctx context.Context) {
	agent.schedule(func() { agent.sendMetricsV1(); agent.sendMetricsV2() }, ctx, agent.config.ReportInterval, "sending metrics")
}

func (agent *Agent) sendMetricsV2() {
	var apiMetrics []protocol.Metrics
	if agent.config.SignKey != "" {
		apiMetrics = agent.metrics.ToAPIWithSign(agent.config.SignKey)
	} else {
		apiMetrics = agent.metrics.ToAPI()
	}

	for i := 0; i < len(apiMetrics); i++ {
		if json, err := json.Marshal(apiMetrics[i]); err != nil {
			agent.logger.Errorf("error in Marshal metric: %s", err)
		} else if r, err := agent.client.Post(fmt.Sprintf("%s/update", agent.url), "application/json", bytes.NewBuffer(json)); err != nil {
			agent.logger.Errorf("error in send metric: %s", err)
		} else {
			r.Body.Close()
		}
	}
}

func (agent *Agent) sendMetricsV1() {
	send := func(url string) {
		if r, err := agent.client.Post(url, "text/plain", nil); err != nil {
			agent.logger.Errorf("error in send metric: %s", err)
		} else {
			r.Body.Close()
		}
	}

	agent.metrics.OperateOverMetricMaps(
		func(key string, v float64) {
			send(fmt.Sprintf("%s/update/gauge/%s/%v", agent.url, key, v))
		},
		func(key string, v int64) {
			send(fmt.Sprintf("%s/update/counter/%s/%v", agent.url, key, v))
		},
	)
}
