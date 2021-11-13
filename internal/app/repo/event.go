package repo

import (
	"context"
	"fmt"
	"sort"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/ozonmp/est-water-api/internal/database"
	"github.com/ozonmp/est-water-api/internal/model"
)

const waterEventTableName = "water_events"

//go:generate mockgen -destination=../../mocks/event_repo_mock.go -package=mocks github.com/ozonmp/est-water-api/internal/app/repo EventRepo
type EventRepo interface {
	Lock(ctx context.Context, n uint64) ([]model.WaterEvent, error)
	Unlock(ctx context.Context, eventIDs []uint64) error
	Add(ctx context.Context, events []model.WaterEvent) error
	Remove(ctx context.Context, eventIDs []uint64) error
}

type eventRepo struct {
	db *sqlx.DB
}

func NewEventRepo(db *sqlx.DB) EventRepo {
	return &eventRepo{db: db}
}

func (er *eventRepo) Lock(ctx context.Context, n uint64) ([]model.WaterEvent, error) {
	subQuery := database.StatementBuilder.
		Select("id").
		From(waterEventTableName).
		OrderBy("id").
		Where(sq.Eq{"status": "unlock"}).
		Limit(n).
		Suffix("FOR NO KEY UPDATE")

	withQuery := subQuery.Prefix("WITH cte AS (").Suffix(")")

	queryText, queryArgs, err := database.StatementBuilder.
		Update(fmt.Sprintf("%s we", waterEventTableName)).
		PrefixExpr(withQuery).
		Set("status", "lock").
		Set("updated_at", time.Now().UTC()).
		Suffix("FROM cte WHERE we.id = cte.id RETURNING we.*").
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "query.ToSql() failed")
	}

	rows, err := er.db.QueryxContext(ctx, queryText, queryArgs...)
	if err != nil {
		return nil, errors.Wrap(err, "db.QueryxContext() failed")
	}
	defer rows.Close()

	waterEvents := make([]model.WaterEvent, 0, n)
	for rows.Next() {
		var waterEvent model.WaterEvent
		if err = rows.StructScan(&waterEvent); err != nil {
			return nil, errors.Wrap(err, "rows.StructScan() failed")
		}
		waterEvents = append(waterEvents, waterEvent)
	}

	sort.Slice(waterEvents, func(i, j int) bool {
		return waterEvents[i].ID < waterEvents[j].ID
	})

	return waterEvents, nil
}

func (er *eventRepo) Unlock(ctx context.Context, eventIDs []uint64) (err error) {
	query := database.StatementBuilder.
		Update(waterEventTableName).
		Set("status", "unlock").
		Set("updated_at", time.Now().UTC()).
		Where(sq.Eq{"id": eventIDs}).
		RunWith(er.db)

	_, err = query.ExecContext(ctx)
	return errors.Wrap(err, "db.ExecContext() failed")
}

func (er *eventRepo) Remove(ctx context.Context, eventIDs []uint64) error {
	query := database.StatementBuilder.
		Delete(waterEventTableName).
		Where(sq.Eq{"id": eventIDs}).
		RunWith(er.db)

	_, err := query.ExecContext(ctx)
	return errors.Wrap(err, "db.ExecContext() failed")
}

func (er *eventRepo) Add(ctx context.Context, events []model.WaterEvent) error {
	query := database.StatementBuilder.
		Insert(waterEventTableName).
		Columns("water_id", "type", "status", "payload", "created_at")

	for _, event := range events {
		query = query.Values(event.WaterId, event.Type, event.Status, event.Entity, event.CreatedAt)
	}

	_, err := query.RunWith(er.db).ExecContext(ctx)
	return errors.Wrap(err, "db.ExecContext() failed")
}