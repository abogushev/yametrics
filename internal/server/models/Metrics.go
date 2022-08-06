package models

const (
	GAUGE   = "gauge"
	COUNTER = "counter"
)

type Metrics struct {
	ID    string
	MType string
	Delta int64
	Value float64
}
