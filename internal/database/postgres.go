package database

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var StatementBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func NewPostgres(ctx context.Context, dsn, driver string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create database connection")

		return nil, err
	}

	 if err = db.PingContext(ctx); err != nil {
		 log.Error().Err(err).Msgf("failed ping the database")
		 return nil, err
	}

	return db, nil
}
