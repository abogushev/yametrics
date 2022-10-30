package storage

import (
	"fmt"
	"yametrics/internal/server/models"
)
//MetricsStorage - интерфейс для абстрагирования работ с хранилищем
type MetricsStorage interface {
	Get(id string, mtype string) (*models.Metrics, error)
	GetAll() ([]models.Metrics, error)
	Update(*models.Metrics) error
	Updates([]models.Metrics) error
	Close()
	Check() error
}

type storageInitError struct {
	err error
}

func (s *storageInitError) Error() string {
	return fmt.Sprintf("error on init storage: %v", s.err)
}

func NewStorageInitError(err error) error {
	return &storageInitError{err}
}
