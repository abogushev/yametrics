package storage

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"
	"yametrics/internal/server/config"
	"yametrics/internal/server/models"

	"go.uber.org/zap"
)

type MetricsStorage interface {
	Get(id string, mtype string) (models.Metrics, bool)
	GetAll() []models.Metrics
	Update(models.Metrics)
}

type metricsStorageImpl struct {
	mutex   sync.Mutex
	metrics map[string]*models.Metrics
	logger  *zap.SugaredLogger
	cfg     *config.MetricsStorageConfig
}

func NewMetricsStorageImpl(
	cfg *config.MetricsStorageConfig,
	logger *zap.SugaredLogger,
	ctx context.Context,
	wg *sync.WaitGroup) (MetricsStorage, error) {
	storage := &metricsStorageImpl{metrics: make(map[string]*models.Metrics), cfg: cfg, logger: logger}
	if cfg.Restore {
		if err := storage.loadMetrics(); err != nil {
			return nil, err
		}
	}
	if cfg.StoreFile != "" {
		go storage.runSaveMetricsJob(ctx, wg)
	}
	return storage, nil
}

func (s *metricsStorageImpl) Get(id string, mtype string) (models.Metrics, bool) {
	if v, ok := s.metrics[id]; ok && v.MType == mtype {
		return *v, true
	} else {
		return models.Metrics{}, false
	}
}

func (s *metricsStorageImpl) GetAll() []models.Metrics {
	m := make([]models.Metrics, len(s.metrics))
	i := 0
	for _, v := range s.metrics {
		m[i] = *v
		i++
	}
	return m
}

func (s *metricsStorageImpl) Update(m models.Metrics) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if v, ok := s.metrics[m.ID]; ok && v.MType == models.COUNTER {
		v.Delta += m.Delta
	} else {
		s.metrics[m.ID] = &m
	}
}

func (s *metricsStorageImpl) runSaveMetricsJob(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(s.cfg.StoreInterval)
	wg.Add(1)
	defer wg.Done()

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

func (s *metricsStorageImpl) saveMetrics() {
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

func (s *metricsStorageImpl) loadMetrics() error {
	s.logger.Info("starting load metrics...")
	if file, err := os.OpenFile(s.cfg.StoreFile, os.O_RDONLY|os.O_CREATE, 0777); err != nil {
		s.logger.Error("error on load metrics", err)
		return err
	} else {
		defer file.Close()
		decoder := json.NewDecoder(file)
		metrics := make(map[string]*models.Metrics, 0)
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
				metrics[m.ID] = &m
			}
		}
	}
}
