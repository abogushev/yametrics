package storage

type CounterStorage interface {
	Update(string, int64)
}

type CounterStorageImpl struct {
	metrics map[string]int64
}

func (s *CounterStorageImpl) Update(name string, value int64) {
	s.metrics[name] += value
}

func NewCounterStorage() CounterStorage {
	return &CounterStorageImpl{map[string]int64{}}
}
