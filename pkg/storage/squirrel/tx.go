package squirrel

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func WithTx(ctx context.Context, db *sqlx.DB, logger *zap.SugaredLogger, op func(*sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	err = op(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logger.With("err", err).Errorf("error rolling back transaction")
		}
		return err
	}

	return tx.Commit()
}
