package storage

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type DbMetricStorage struct {
	url string
	ctx context.Context
}

func NewDbMetricStorage(url string, ctx context.Context) *DbMetricStorage {
	return &DbMetricStorage{url, ctx}
}

func (db *DbMetricStorage) Ping() error {
	conn, err := pgx.Connect(context.Background(), db.url)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	err = conn.Ping(db.ctx)
	return err
}
