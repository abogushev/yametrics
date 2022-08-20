package models

const (
	GAUGE   = "gauge"
	COUNTER = "counter"
)

type Metrics struct {
	ID    string   `db:"id"`
	MType string   `db:"mtype"`
	Delta *int64   `db:"delta"`
	Value *float64 `db:"value"`
}
