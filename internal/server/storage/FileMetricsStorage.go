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

// fileMetricsStorage - файловое хранилище метрик
// умеет периодически сохранять даные в файл и вычитывает его при старте
type fileMetricsStorage struct {
	mutex   sync.Mutex
	metrics map[string]*models.Metrics
	logger  *zap.SugaredLogger
	cfg     *config.MetricsStorageConfig
}

func NewFileMetricsStorage(
	cfg *config.MetricsStorageConfig,
	logger *zap.SugaredLogger,
	ctx context.Context) (MetricsStorage, error) {
	storage := &fileMetricsStorage{metrics: make(map[string]*models.Metrics), cfg: cfg, logger: logger}
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

func (s *fileMetricsStorage) Updates(metrics []models.Metrics) error {
	for i := 0; i < len(metrics); i++ {
		s.Update(&metrics[i])
	}
	return nil
}

func (s *fileMetricsStorage) Get(id string, mtype string) (*models.Metrics, error) {
	if metric, ok := s.metrics[id]; ok && metric.MType == mtype {
		v := *metric
		return &v, nil
	} else {
		return nil, nil
	}
}

func (s *fileMetricsStorage) GetAll() ([]models.Metrics, error) {
	m := make([]models.Metrics, len(s.metrics))
	i := 0
	for _, v := range s.metrics {
		m[i] = *v
		i++
	}
	return m, nil
}

func (s *fileMetricsStorage) Update(m *models.Metrics) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if v, ok := s.metrics[m.ID]; ok && v.MType == models.COUNTER {
		*v.Delta += *m.Delta
	} else {
		s.metrics[m.ID] = m
	}
	return nil
}

func (s *fileMetricsStorage) Close() {
	s.saveMetrics()
}

func (s *fileMetricsStorage) Check() error {
	return nil
}

func (s *fileMetricsStorage) runSaveMetricsJob(ctx context.Context) {
	ticker := time.NewTicker(s.cfg.StoreInterval)

	for {
		select {
		case <-ticker.C:
			s.saveMetrics()

		case <-ctx.Done():
			ticker.Stop()
			s.logger.Info("stop runSaveMetricsJob")
			return
		}
	}
}

func (s *fileMetricsStorage) saveMetrics() {
	s.logger.Info("starting save metrics...")
	if file, err := os.OpenFile(s.cfg.StoreFile, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0777); err != nil {
		s.logger.Errorf("error on save metrics: %w", err)

	} else {
		defer file.Close()
		encoder := json.NewEncoder(file)
		for _, m := range s.metrics {
			encoder.Encode(m)
		}
		s.logger.Info("metrics saved")
	}
}

func (s *fileMetricsStorage) loadMetrics() error {
	s.logger.Info("starting load metrics...")
	if file, err := os.OpenFile(s.cfg.StoreFile, os.O_RDONLY|os.O_CREATE, 0777); err != nil {
		s.logger.Errorf("error on load metrics: %w", err)
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
				s.logger.Errorf("error on read from file: %w", err)
				return err
			} else {
				metrics[m.ID] = &m
			}
		}
	}
}
