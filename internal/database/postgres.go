package database

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/ozonmp/est-water-api/internal/logger"
	"github.com/pkg/errors"
)

var StatementBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func NewPostgres(ctx context.Context, dsn, driver string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		logger.ErrorKV(ctx, "failed to create database connection",
			"err", errors.Wrap(err, "sqlx.Open() failed"),
		)

		return nil, err
	}

	 if err = db.PingContext(ctx); err != nil {
		 logger.ErrorKV(ctx, "failed ping the database",
			 "err", errors.Wrap(err, "db.PingContext() failed"),
		 )
		 return nil, err
	}

	return db, nil
}
