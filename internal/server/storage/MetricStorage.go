package storage

import "yametrics/internal/server/models"

type MetricsStorage interface {
	Get(id string, mtype string) (*models.Metrics, error)
	GetAll() ([]models.Metrics, error)
	Update(*models.Metrics) error
	Updates([]models.Metrics) error
	Close()
	Check() error
}
