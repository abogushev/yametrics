package storage

type GuageStorage interface {
	Update(string, float64)
}

type GuageStorageImpl struct {
	metrics map[string]float64
}

func (s *GuageStorageImpl) Update(name string, value float64) {
	s.metrics[name] = value
}

func NewGuageStorage() GuageStorage {
	return &GuageStorageImpl{map[string]float64{}}
}
