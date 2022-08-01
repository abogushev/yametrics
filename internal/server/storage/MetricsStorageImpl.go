package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
	"yametrics/internal/server/models"

	"go.uber.org/zap"
)

type MetricsStorage interface {
	Get(string) (models.Metrics, bool)
	GetAll() []models.Metrics
	Update(models.Metrics)
}

type MetricsStorageImpl struct {
	mutex   sync.Mutex
	metrics map[string]models.Metrics
	logger  *zap.SugaredLogger
	cfg     *models.MetricsStorageConfig
}

func NewMetricsStorageImpl(cfg *models.MetricsStorageConfig, logger *zap.SugaredLogger, ctx context.Context) (MetricsStorage, error) {
	storage := &MetricsStorageImpl{metrics: map[string]models.Metrics{}, cfg: cfg, logger: logger}
	if cfg.Restore {
		if err := storage.loadMetrics(); err != nil {
			return nil, err
		}
	}
	if cfg.StoreFile != "" {
		go storage.runSaveMetricsJob(ctx)
	}
	return storage, nil
}

func (s *MetricsStorageImpl) Get(name string) (v models.Metrics, ok bool) {
	v, ok = s.metrics[name]
	fmt.Printf("get mtrc - key: %v, value: %v, delta: %v, all: %v\n", name, v.Value, v.Delta, s.metrics)
	return
}

func (s *MetricsStorageImpl) GetAll() []models.Metrics {
	m := make([]models.Metrics, len(s.metrics))
	i := 0
	for _, v := range s.metrics {
		m[i] = v
		i++
	}
	return m
}

func (s *MetricsStorageImpl) Update(m models.Metrics) {
	s.mutex.Lock()
	s.metrics[m.ID] = m
	fmt.Printf("update mtrcs - key: %v, value: %v, delta: %v, all: %v\n", m.ID, m.Value, m.Delta, s.metrics)
	s.mutex.Unlock()
}

func (s *MetricsStorageImpl) runSaveMetricsJob(ctx context.Context) {
	ticker := time.NewTicker(s.cfg.StoreInterval)
	for {
		select {
		case <-ticker.C:
			s.saveMetrics()

		case <-ctx.Done():
			ticker.Stop()
			s.saveMetrics()
			s.logger.Info("stop runSaveMetricsJob")
			return
		}
	}
}

func (s *MetricsStorageImpl) saveMetrics() {
	s.logger.Info("starting save metrics...")
	if file, err := os.OpenFile(s.cfg.StoreFile, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0777); err != nil {
		s.logger.Error("error on save metrics", err)

	} else {
		defer file.Close()
		encoder := json.NewEncoder(file)
		for _, m := range s.metrics {
			encoder.Encode(m)
		}
		s.logger.Info("metrics saved")
	}
}

func (s *MetricsStorageImpl) loadMetrics() error {
	s.logger.Info("starting load metrics...")
	if file, err := os.OpenFile(s.cfg.StoreFile, os.O_RDONLY|os.O_CREATE, 0777); err != nil {
		s.logger.Error("error on load metrics", err)
		return err
	} else {
		defer file.Close()
		decoder := json.NewDecoder(file)
		metrics := make(map[string]models.Metrics, 0)
		for {
			var m models.Metrics
			if err := decoder.Decode(&m); err == io.EOF {
				s.metrics = metrics
				s.logger.Info("load metrics completed")
				return nil
			} else if err != nil {
				s.logger.Error("error on read from file", err)
				return err
			} else {
				metrics[m.ID] = m
			}
		}
	}
}
