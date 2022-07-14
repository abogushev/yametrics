package storage

type CounterStorage interface {
	Get(string) (int64, bool)
	GetAll() map[string]int64
	Update(string, int64)
}

type CounterStorageImpl struct {
	metrics map[string]int64
}

func (s *CounterStorageImpl) Update(name string, value int64) {
	s.metrics[name] += value
}

func (s *CounterStorageImpl) Get(name string) (v int64, ok bool) {
	v, ok = s.metrics[name]
	return
}

func (s *CounterStorageImpl) GetAll() map[string]int64 {
	return s.metrics
}

func NewCounterStorage() CounterStorage {
	return &CounterStorageImpl{map[string]int64{}}
}
