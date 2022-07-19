package storage

import "sync"

type CounterStorage interface {
	Get(string) (int64, bool)
	GetAll() map[string]int64
	Update(string, int64)
}

type CounterStorageImpl struct {
	mutex   sync.Mutex
	metrics map[string]int64
}

func (s *CounterStorageImpl) Update(name string, value int64) {
	s.mutex.Lock()

	s.metrics[name] += value
	s.mutex.Unlock()
}

func (s *CounterStorageImpl) Get(name string) (v int64, ok bool) {
	v, ok = s.metrics[name]
	return
}

func (s *CounterStorageImpl) GetAll() map[string]int64 {
	return s.metrics
}

func NewCounterStorage() CounterStorage {
	return &CounterStorageImpl{metrics: map[string]int64{}}
}
