package storage

type GuageStorage interface {
	Get(string) (float64, bool)
	GetAll() map[string]float64
	Update(string, float64)
}

type GuageStorageImpl struct {
	metrics map[string]float64
}

func (s *GuageStorageImpl) Update(name string, value float64) {
	s.metrics[name] = value
}

func (s *GuageStorageImpl) Get(name string) (v float64, ok bool) {
	v, ok = s.metrics[name]
	return
}

func (s *GuageStorageImpl) GetAll() map[string]float64 {
	return s.metrics
}

func NewGuageStorage() GuageStorage {
	return &GuageStorageImpl{map[string]float64{}}
}
