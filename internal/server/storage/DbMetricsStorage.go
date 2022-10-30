package storage

import (
	"context"
	"database/sql"
	"errors"
	"yametrics/internal/server/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	createTableIfNeedSQL = `create table if not exists metrics(
		id varchar not null primary key,
		mtype varchar not null,
		delta bigint,
		value double precision)`

	upInsertSQL = `insert into metrics(
		id,
		mtype,
		delta,
		value)

		values(
		:id,
		:mtype,
		:delta,
		:value)
		
		on conflict(id) do update set 
		mtype = :mtype, 
		delta = case when metrics.mtype = 'counter' then metrics.delta + :delta end,
		value = case when metrics.mtype = 'gauge' then CAST(:value AS DOUBLE PRECISION) end`
	getSQL    = `select id, mtype, delta, value from metrics where id = $1 and mtype = $2`
	getAllSQL = `select id, mtype, delta, value from metrics`
)
//dbMetricStorage - сервис по работе с бд
type dbMetricStorage struct {
	url    string
	ctx    context.Context
	xdb    *sqlx.DB
	logger *zap.SugaredLogger
}

func NewDBMetricStorage(url string, ctx context.Context, logger *zap.SugaredLogger) (MetricsStorage, error) {
	logger.Infow("start init dbstorage ...")
	xdb, err := sqlx.Connect("postgres", url)
	if err != nil {
		logger.Errorf("error on connect to db: %v", err)
		return nil, NewStorageInitError(err)
	}

	storage := &dbMetricStorage{url, ctx, xdb, logger}
	if err := storage.initDB(); err != nil {
		logger.Errorf("error on connect to init db: %v", err)
		return nil, NewStorageInitError(err)
	}
	logger.Info("dbstorage initialized successfully")
	return storage, nil
}

var upInsertStmt *sqlx.NamedStmt

func (db *dbMetricStorage) initDB() error {
	_, err := db.xdb.ExecContext(db.ctx, createTableIfNeedSQL)

	if err != nil {
		return err
	}

	upInsertStmt, err = db.xdb.PrepareNamed(upInsertSQL)

	return err
}

func (db *dbMetricStorage) Get(id string, mtype string) (*models.Metrics, error) {
	metric := models.Metrics{}
	err := db.xdb.GetContext(db.ctx, &metric, getSQL, id, mtype)
	if err == nil {
		return &metric, nil
	} else if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else {
		return nil, err
	}
}

func (db *dbMetricStorage) GetAll() ([]models.Metrics, error) {
	metrics := []models.Metrics{}
	if err := db.xdb.SelectContext(db.ctx, &metrics, getAllSQL); err != nil {
		return nil, err
	}
	return metrics, nil

}

func (db *dbMetricStorage) Update(m *models.Metrics) error {
	if _, err := upInsertStmt.ExecContext(db.ctx, m); err != nil {
		return err
	}
	return nil
}

func (db *dbMetricStorage) Updates(mtrcs []models.Metrics) error {
	tx, err := db.xdb.Beginx()
	if err != nil {
		return err
	}

	for i := 0; i < len(mtrcs); i++ {
		if _, err := upInsertStmt.ExecContext(db.ctx, &mtrcs[i]); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}
	return tx.Commit()
}

func (db *dbMetricStorage) Check() error {
	return db.xdb.PingContext(db.ctx)
}

func (db *dbMetricStorage) Close() {
	db.logger.Info("closing dbMetricStorage...")
	db.xdb.Close()
	db.logger.Info("dbMetricStorage closed")
}
