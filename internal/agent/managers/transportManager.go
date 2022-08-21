package managers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"

	"yametrics/internal/agent/config"
	"yametrics/internal/agent/models/storage"
	"yametrics/internal/agent/utils"
	"yametrics/internal/protocol"

	"go.uber.org/zap"
)

type transportManager struct {
	url     string
	client  http.Client
	logger  *zap.SugaredLogger
	config  *config.AgentConfig
	metrics *storage.Metrics
	rwMutex sync.RWMutex
	once    sync.Once
}

func NewTransportManager(l *zap.SugaredLogger, config *config.AgentConfig) *transportManager {
	return &transportManager{
		url:     "http://" + config.Address,
		client:  http.Client{},
		logger:  l,
		metrics: &storage.Metrics{MemStats: &runtime.MemStats{}, PollCount: 0, RandomValue: 0.0},
		config:  config,
	}
}

func (m *transportManager) RunAsync(notifyCh <-chan storage.Metrics, ctx context.Context, wg *sync.WaitGroup) {
	m.once.Do(func() {
		wg.Add(2)
		go m.sendMetricsWithInterval(ctx, wg)
		go m.subscribeOnUpdates(notifyCh, ctx, wg)
	})
}

func (m *transportManager) subscribeOnUpdates(notifyCh <-chan storage.Metrics, ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case updated := <-notifyCh:
			m.logger.Info("update metrics from notifyCh")
			m.rwMutex.Lock()
			m.metrics = &updated
			m.rwMutex.Unlock()

		case <-ctx.Done():
			m.logger.Info("transportManager shutdown")
			wg.Done()
			return
		}
	}
}

func (m *transportManager) sendMetricsWithInterval(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	utils.Schedule(func() {
		m.rwMutex.RLock()
		defer m.rwMutex.RUnlock()

		m.sendMetricsV1()
		m.sendMetricsV2()
		m.sendMultipleMetricsV2()
	},
		ctx,
		m.config.ReportInterval,
		"sending metrics",
		m.logger)
}

func (m *transportManager) sendMultipleMetricsV2() {
	apiMetrics := m.metrics.ToAPI()
	if json, err := json.Marshal(apiMetrics); err != nil {
		m.logger.Errorf("error on  Marshal metric: %w", err)
	} else {
		if r, err := m.client.Post(fmt.Sprintf("%s/updates", m.url), "application/json", bytes.NewBuffer(json)); err != nil {
			m.logger.Errorf("error in send metric: %w", err)
		} else {
			r.Body.Close()
		}
	}
}

func (m *transportManager) sendMetricsV2() {
	var apiMetrics []protocol.Metrics
	if m.config.SignKey != "" {
		apiMetrics = m.metrics.ToAPIWithSign(m.config.SignKey)
	} else {
		apiMetrics = m.metrics.ToAPI()
	}

	for i := 0; i < len(apiMetrics); i++ {
		if json, err := json.Marshal(apiMetrics[i]); err != nil {
			m.logger.Errorf("error in Marshal metric: %w", err)
		} else if r, err := m.client.Post(fmt.Sprintf("%s/update", m.url), "application/json", bytes.NewBuffer(json)); err != nil {
			m.logger.Errorf("error in send metric: %w", err)
		} else {
			r.Body.Close()
		}
	}
}

func (m *transportManager) sendMetricsV1() {
	send := func(url string) {
		if r, err := m.client.Post(url, "text/plain", nil); err != nil {
			m.logger.Errorf("error in send metric: %w", err)
		} else {
			r.Body.Close()
		}
	}

	m.metrics.OperateOverMetricMaps(
		func(key string, v float64) {
			send(fmt.Sprintf("%s/update/gauge/%s/%v", m.url, key, v))
		},
		func(key string, v int64) {
			send(fmt.Sprintf("%s/update/counter/%s/%v", m.url, key, v))
		},
	)
}
