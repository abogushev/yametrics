package storage

import (
	"context"
	"yametrics/internal/server/models"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBMetricStorage struct {
	url    string
	ctx    context.Context
	dbpool *pgxpool.Pool
}

func NewDBMetricStorage(url string, ctx context.Context) (MetricsStorage, error) {
	dbpool, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	storage := &DBMetricStorage{url, ctx, dbpool}
	if err := storage.initDB(); err != nil {
		return nil, err
	}
	return storage, nil
}

func (db *DBMetricStorage) initDB() error {
	_, err := db.dbpool.Exec(db.ctx, "create table if not exists metrics(id varchar not null primary key, mtype varchar not null, delta bigint, value double precision)")
	return err
}

func (db *DBMetricStorage) Get(id string, mtype string) (*models.Metrics, error) {
	metric := models.Metrics{}
	row := db.dbpool.QueryRow(db.ctx, "select id, mtype, delta, value from metrics where id = $1 and mtype = $2", id, mtype)
	err := row.Scan(&metric.ID, &metric.MType, metric.Delta, metric.Value)
	if err == nil {
		return &metric, nil
	} else if err == pgx.ErrNoRows {
		return nil, nil
	} else {
		return nil, err
	}
}

func (db *DBMetricStorage) GetAll() ([]models.Metrics, error) {
	metrics := make([]models.Metrics, 0)
	rows, err := db.dbpool.Query(db.ctx, "select id, mtype, delta, value from metrics")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		metric := models.Metrics{}
		err := rows.Scan(&metric.ID, &metric.MType, metric.Delta, metric.Value)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return metrics, nil
}

func (db *DBMetricStorage) Update(m *models.Metrics) error {
	if stored, err := db.Get(m.ID, m.MType); err != nil {
		return err
	} else if stored != nil && stored.MType == models.COUNTER {
		*m.Delta += *stored.Delta
		if _, err = db.dbpool.Exec(db.ctx, "update metrics set mtype = $1, delta = $2, value = $3 where id = $4", m.MType, m.Delta, m.Value, m.ID); err != nil {
			return err
		}
	} else if _, err = db.dbpool.Exec(db.ctx, "insert into metrics(id, mtype, delta, value) values($1,$2,$3,$4) on conflict(id) do update set mtype = $2, delta = $3, value = $4", m.ID, m.MType, m.Delta, m.Value); err != nil {
		return err
	}
	return nil
}

func (db *DBMetricStorage) Check() error {
	return db.dbpool.Ping(db.ctx)
}

func (db *DBMetricStorage) Close() {
	db.dbpool.Close()
}
