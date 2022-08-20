package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"yametrics/internal/agent/config"
	"yametrics/internal/agent/models/storage"
	"yametrics/internal/protocol"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
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
	wg := &sync.WaitGroup{}

	go agent.updateMetricsWithInterval(ctx, wg)
	go agent.collectAdditional(ctx, wg)

	go agent.sendMetricsWithInterval(ctx, wg)

	agent.logger.Info("agent started")
	wg.Wait()
	agent.logger.Info("agent stoped")
}

func (agent *Agent) collectAdditional(ctx context.Context, wg *sync.WaitGroup) {
	v, _ := mem.VirtualMemoryWithContext(ctx)
	agent.metrics.TotalMemory = float64(v.Total)
	agent.metrics.FreeMemory = float64(v.Free)
	cpuUsage, _ := cpu.PercentWithContext(ctx, 0, false)
	agent.metrics.CPUutilization1 = cpuUsage[0]
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

func (agent *Agent) updateMetricsWithInterval(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
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

func (agent *Agent) sendMetricsWithInterval(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	agent.schedule(func() { agent.sendMetricsV1(); agent.sendMetricsV2(); agent.sendMultipleMetricsV2() }, ctx, agent.config.ReportInterval, "sending metrics")
}

func (agent *Agent) sendMultipleMetricsV2() {
	apiMetrics := agent.metrics.ToAPI()
	if json, err := json.Marshal(apiMetrics); err != nil {
		agent.logger.Errorf("error on  Marshal metric: %w", err)
	} else {
		if r, err := agent.client.Post(fmt.Sprintf("%s/updates", agent.url), "application/json", bytes.NewBuffer(json)); err != nil {
			agent.logger.Errorf("error in send metric: %w", err)
		} else {
			r.Body.Close()
		}
	}
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
			agent.logger.Errorf("error in Marshal metric: %w", err)
		} else if r, err := agent.client.Post(fmt.Sprintf("%s/update", agent.url), "application/json", bytes.NewBuffer(json)); err != nil {
			agent.logger.Errorf("error in send metric: %w", err)
		} else {
			r.Body.Close()
		}
	}
}

func (agent *Agent) sendMetricsV1() {
	send := func(url string) {
		if r, err := agent.client.Post(url, "text/plain", nil); err != nil {
			agent.logger.Errorf("error in send metric: %w", err)
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
