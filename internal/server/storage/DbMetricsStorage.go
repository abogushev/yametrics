package storage

import (
	"context"
	"database/sql"
	"errors"
	"yametrics/internal/server/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBMetricStorage struct {
	url string
	ctx context.Context
	xdb *sqlx.DB
}

func NewDBMetricStorage(url string, ctx context.Context) (MetricsStorage, error) {
	xdb, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}
	storage := &DBMetricStorage{url, ctx, xdb}
	if err := storage.initDB(); err != nil {
		return nil, err
	}
	return storage, nil
}

var insertStmt *sqlx.NamedStmt

func (db *DBMetricStorage) initDB() error {
	_, err := db.xdb.ExecContext(db.ctx, "create table if not exists metrics(id varchar not null primary key, mtype varchar not null, delta bigint, value double precision)")

	if err != nil {
		return err
	}

	insertStmt, err = db.xdb.PrepareNamed("insert into metrics(id, mtype, delta, value) values(:id, :mtype, :delta, :value) on conflict(id) do update set mtype = :mtype, delta = :delta, value = :value")

	return err
}

func (db *DBMetricStorage) Get(id string, mtype string) (*models.Metrics, error) {
	metric := models.Metrics{}
	err := db.xdb.GetContext(db.ctx, &metric, "select id, mtype, delta, value from metrics where id = $1 and mtype = $2", id, mtype)
	if err == nil {
		return &metric, nil
	} else if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else {
		return nil, err
	}
}

func (db *DBMetricStorage) GetAll() ([]models.Metrics, error) {
	metrics := make([]models.Metrics, 0)
	if err := db.xdb.SelectContext(db.ctx, metrics, "select id, mtype, delta, value from metrics"); err != nil {
		return nil, err
	} else {
		return metrics, nil
	}
}

func (db *DBMetricStorage) Update(m *models.Metrics) error {
	if stored, err := db.Get(m.ID, m.MType); err != nil {
		return err
	} else if stored != nil && stored.MType == models.COUNTER {
		*m.Delta += *stored.Delta
		if _, err = db.xdb.ExecContext(db.ctx, "update metrics set mtype = $1, delta = $2, value = $3 where id = $4", m.MType, m.Delta, m.Value, m.ID); err != nil {
			return err
		}
	} else if _, err = insertStmt.ExecContext(db.ctx, m); err != nil {
		return err
	}
	return nil
}

func (db *DBMetricStorage) Updates(mtrcs []models.Metrics) error {
	tx, err := db.xdb.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i := 0; i < len(mtrcs); i++ {
		if _, err := insertStmt.ExecContext(db.ctx, &mtrcs[i]); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}
	return tx.Commit()
}

func (db *DBMetricStorage) Check() error {
	return db.xdb.PingContext(db.ctx)
}

func (db *DBMetricStorage) Close() {
	db.xdb.Close()
}
