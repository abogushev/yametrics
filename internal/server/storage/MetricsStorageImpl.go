package storage

import (
	"sync"
	"yametrics/internal/server/models"
)

type MetricsStorage interface {
	Get(string) (models.Metrics, bool)
	GetAll() []models.Metrics
	Update(models.Metrics)
}

type MetricsStorageImpl struct {
	mutex   sync.Mutex
	metrics map[string]models.Metrics
}

func NewMetricsStorageImpl() MetricsStorage {
	return &MetricsStorageImpl{metrics: map[string]models.Metrics{}}
}

func (s *MetricsStorageImpl) Get(name string) (v models.Metrics, ok bool) {
	v, ok = s.metrics[name]
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
	s.mutex.Unlock()
}
