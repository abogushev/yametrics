package managers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"yametrics/internal/agent/config"
	"yametrics/internal/agent/models/storage"
	"yametrics/internal/agent/utils"
	"yametrics/internal/protocol"

	"go.uber.org/zap"
)

// TransportManager менеджер отправки данных на сервер
type TransportManager struct {
	url     string
	client  http.Client
	logger  *zap.SugaredLogger
	config  *config.AgentConfig
	metrics *storage.Metrics
	rwMutex sync.RWMutex
	once    sync.Once
}

// NewTransportManager - создание менеджера отправки метрик.
// для запуска менеждера необходимо вызвать RunAsync.
func NewTransportManager(l *zap.SugaredLogger, config *config.AgentConfig) *TransportManager {
	return &TransportManager{
		url:     "http://" + config.Address,
		client:  http.Client{},
		logger:  l,
		metrics: storage.NewMetrics(),
		config:  config,
	}
}

// RunAsync - запуск менеджера: обновленные данные будут приходить из notifyCh канала
// здесь стартуют рутины по прослушиванию обновленных данных и по отправке данных на сервер
func (m *TransportManager) RunAsync(notifyCh <-chan storage.Metrics, ctx context.Context, wg *sync.WaitGroup) {
	m.once.Do(func() {
		wg.Add(2)
		go m.sendMetricsWithInterval(ctx, wg)
		go m.subscribeOnUpdates(notifyCh, ctx, wg)
	})
}

func (m *TransportManager) subscribeOnUpdates(notifyCh <-chan storage.Metrics, ctx context.Context, wg *sync.WaitGroup) {
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

func (m *TransportManager) sendMetricsWithInterval(ctx context.Context, wg *sync.WaitGroup) {
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
func (m *TransportManager) sendMultipleMetricsV2() {
	apiMetrics := m.metrics.ToAPI()
	if marshal, err := json.Marshal(apiMetrics); err != nil {
		m.logger.Errorf("error on  Marshal metric: %v", err)
	} else {
		if r, err := m.client.Post(fmt.Sprintf("%s/updates", m.url), "application/json", bytes.NewBuffer(marshal)); err != nil {
			m.logger.Errorf("error in send metric: %v", err)
		} else {
			if err := r.Body.Close(); err != nil {
				m.logger.Errorf("error in close body %v", err)
			}
		}
	}
}

func (m *TransportManager) sendMetricsV2() {
	var apiMetrics []protocol.Metrics
	if m.config.SignKey != "" {
		apiMetrics = m.metrics.ToAPIWithSign(m.config.SignKey)
	} else {
		apiMetrics = m.metrics.ToAPI()
	}

	for i := 0; i < len(apiMetrics); i++ {
		if marshal, err := json.Marshal(apiMetrics[i]); err != nil {
			m.logger.Errorf("error in Marshal metric: %v", err)
		} else if r, err := m.client.Post(fmt.Sprintf("%s/update", m.url), "application/json", bytes.NewBuffer(marshal)); err != nil {
			m.logger.Errorf("error in send metric: %v", err)
		} else {
			if err := r.Body.Close(); err != nil {
				m.logger.Errorf("error in close body %v", err)
			}
		}
	}
}

func (m *TransportManager) sendMetricsV1() {
	send := func(url string) {
		if r, err := m.client.Post(url, "text/plain", nil); err != nil {
			m.logger.Errorf("error in send metric: %v", err)
		} else {
			if err := r.Body.Close(); err != nil {
				m.logger.Errorf("error in close body %v", err)
			}
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
