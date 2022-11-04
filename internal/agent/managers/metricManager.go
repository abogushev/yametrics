package managers

import (
	"context"
	"math/rand"
	"runtime"
	"sync"
	"yametrics/internal/agent/config"
	"yametrics/internal/agent/models/storage"
	"yametrics/internal/agent/utils"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"go.uber.org/zap"
)

type MetricWorker int

const (
	General MetricWorker = iota
	Additional
)

type metricManager struct {
	logger  *zap.SugaredLogger
	metrics *storage.Metrics
	config  *config.AgentConfig
	syncCh  chan MetricWorker
	// NotifyCh - в этот канал будут отдаваться собранные данные.
	NotifyCh      chan storage.Metrics
	syncWorkersMu sync.Mutex
	once          sync.Once
}

// NewMetricManager - создание менеджера сбора метрик.
//
// для запуска менеждера необходимо вызвать RunAsync.
func NewMetricManager(
	logger *zap.SugaredLogger,
	config *config.AgentConfig) *metricManager {

	return &metricManager{
		logger:   logger,
		metrics:  &storage.Metrics{MemStats: &runtime.MemStats{}, PollCount: 0, RandomValue: 0.0},
		config:   config,
		syncCh:   make(chan MetricWorker),
		NotifyCh: make(chan storage.Metrics),
	}
}

// RunAsync - запуск менеджера: стартуют рутины по сбору и отправке метрик
func (m *metricManager) RunAsync(ctx context.Context, wg *sync.WaitGroup) {
	m.once.Do(func() {
		wg.Add(3)
		go m.updateMetricsWithInterval(ctx, wg)
		go m.updateAdditionalMetricsWithInterval(ctx, wg)
		go m.syncUpdates(ctx, wg)
	})
}

// метод для синхронизации джоб по сбору метрик.
func (m *metricManager) syncUpdates(ctx context.Context, wg *sync.WaitGroup) {
	unicJob := make(map[MetricWorker]struct{})

	for {
		select {
		case w := <-m.syncCh:
			m.logger.Debug("metrics updated by %v", w)
			if v, ok := unicJob[w]; !ok {
				unicJob[w] = v
			}
			if len(unicJob) == 2 {
				unicJob = make(map[MetricWorker]struct{})
				m.logger.Debug("sending updated metrics to channel")
				m.NotifyCh <- *m.metrics
			}

		case <-ctx.Done():
			m.logger.Info("metricsManager shutdown")
			close(m.NotifyCh)
			wg.Done()
			return
		}
	}
}

func (m *metricManager) updateAdditionalMetricsWithInterval(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	utils.Schedule(
		func() {
			m.syncWorkersMu.Lock()
			defer m.syncWorkersMu.Unlock()

			v, _ := mem.VirtualMemoryWithContext(ctx)
			m.metrics.TotalMemory = float64(v.Total)
			m.metrics.FreeMemory = float64(v.Free)
			cpuUsage, err := cpu.PercentWithContext(ctx, 0, false)
			if err != nil {
				m.logger.Errorf("err in updateAdditionalMetricsWithInterval %v", err)
			}
			m.metrics.CPUutilization1 = cpuUsage[0]

			m.syncCh <- Additional
		},
		ctx,
		m.config.PollInterval,
		"collecting additional metrics",
		m.logger)
}

func (m *metricManager) updateMetricsWithInterval(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	utils.Schedule(
		func() {
			m.syncWorkersMu.Lock()

			runtime.ReadMemStats(m.metrics.MemStats)
			m.metrics.PollCount++
			m.metrics.RandomValue = rand.Float64()

			m.syncWorkersMu.Unlock()

			m.syncCh <- General
		},
		ctx,
		m.config.PollInterval,
		"collecting metrics",
		m.logger)
}
