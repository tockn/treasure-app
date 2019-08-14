package rdb

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	datasource string
}

func NewDB(datasource string) *DB {
	return &DB{
		datasource: datasource,
	}
}

func (db *DB) Open() (*sqlx.DB, error) {
	return sqlx.Open("mysql", db.datasource)
}

var txKey = struct{}{}

func GetTx(ctx context.Context) (Tx, bool) {
	tx, ok := ctx.Value(&txKey).(Tx)
	return tx, ok
}

type Tx interface {
	Prepare(query string) (*sql.Stmt, error)
}

type txHandler struct {
	db *sqlx.DB
}

func (r *txHandler) TxHandler(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// sqlx.Txで平気？ txHandler.Txにしなくて良い?
	ctx = context.WithValue(ctx, &txKey, tx)
	if err := fn(ctx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
