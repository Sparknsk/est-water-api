package database

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var StatementBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func NewPostgres(ctx context.Context, dsn, driver string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, errors.Wrap(err, "sqlx.Open() failed")
	}

	 if err = db.PingContext(ctx); err != nil {
		 return nil, errors.Wrap(err, "db.PingContext() failed")
	}

	return db, nil
}
