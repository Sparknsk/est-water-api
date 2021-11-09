package repo

import (
	"context"
	"fmt"
	"sort"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/ozonmp/est-water-api/internal/database"
	"github.com/ozonmp/est-water-api/internal/model"
)

const waterEventTableName = "water_events"

//go:generate mockgen -destination=../../mocks/repo_mock.go -package=mocks github.com/ozonmp/est-water-api/internal/app/repo EventRepo
type EventRepo interface {
	Lock(ctx context.Context, n uint64) ([]model.WaterEvent, error)
	Unlock(ctx context.Context, eventIDs []uint64) error

	Add(ctx context.Context, event *model.WaterEvent) error
	Remove(ctx context.Context, eventIDs []uint64) error
}

type eventRepo struct {
	db *sqlx.DB
}

func NewEventRepo(db *sqlx.DB) EventRepo {
	return &eventRepo{db: db}
}

func (er *eventRepo) Lock(ctx context.Context, n uint64) ([]model.WaterEvent, error) {
	//TODO: потестировать потом NOWAIT, может будет быстрее обработать ошибку в consumer и пойти снова в базу за событиями
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
		Suffix("FROM cte WHERE we.id = cte.id RETURNING we.*").
		ToSql()

	rows, err := er.db.QueryxContext(ctx, queryText, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	waterEvents := make([]model.WaterEvent, 0, n)
	for rows.Next() {
		var waterEvent model.WaterEvent
		if err = rows.StructScan(&waterEvent); err != nil {
			return nil, err
		}
		waterEvents = append(waterEvents, waterEvent)
	}

	sort.Slice(waterEvents, func(i, j int) bool {
		return waterEvents[i].ID < waterEvents[j].ID
	})

	return waterEvents, nil
}

// Lock TODO: блокировка с воркшопа, возможно получиться потестировать на производительность после следующих воркшопов, пока оставлю (+ не уверен, что верно использую ее в цикле, надо будет уточнить)
/*func (er *eventRepo) Lock(ctx context.Context, n uint64) ([]model.WaterEvent, error) {
	tx, err := er.db.Beginx()
	defer tx.Commit()
	for {
		if err != nil {
			return nil, err
		}

		var isAcquired bool
		lockQueryText := "SELECT pg_try_advisory_xact_lock(100)"
		err = tx.GetContext(ctx, &isAcquired, lockQueryText)
		if err != nil {
			return nil, err
		}

		if isAcquired == true {
			break
		}
	}

	query := "WITH cte AS (SELECT id FROM water_events WHERE status = 'unlock' ORDER BY id LIMIT $1) UPDATE water_events we SET status = 'lock' FROM cte WHERE we.id = cte.id RETURNING we.*;"
	rows, err := tx.QueryxContext(ctx, query, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var waterEvents []model.WaterEvent
	for rows.Next() {
		var waterEvent model.WaterEvent
		if err = rows.StructScan(&waterEvent); err != nil {
			return nil, err
		}
		waterEvents = append(waterEvents, waterEvent)
	}

	sort.Slice(waterEvents, func(i, j int) bool {
		return waterEvents[i].ID < waterEvents[j].ID
	})

	return waterEvents, nil
}*/

func (er *eventRepo) Unlock(ctx context.Context, eventIDs []uint64) (err error) {
	query := database.StatementBuilder.
		Update(waterEventTableName).
		Set("status", "unlock").
		Set("updated_at", time.Now().UTC()).
		Where(sq.Eq{"id": eventIDs})

	queryText, queryArgs, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = er.db.ExecContext(ctx, queryText, queryArgs...)

	return err
}

func (er *eventRepo) Remove(ctx context.Context, eventIDs []uint64) error {
	query := database.StatementBuilder.
		Delete(waterEventTableName).
		Where(sq.Eq{"id": eventIDs})

	queryText, queryArgs, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = er.db.ExecContext(ctx, queryText, queryArgs...)

	return err
}

func (er *eventRepo) Add(ctx context.Context, event *model.WaterEvent) error {
	query := database.StatementBuilder.
		Insert(waterEventTableName).
		Columns("water_id", "type", "status", "payload", "created_at").
		Values(event.WaterId, event.Type, event.Status, event.Entity, event.CreatedAt).
		Suffix("RETURNING id").
		RunWith(er.db)

	if err := query.QueryRowContext(ctx).Scan(&event.ID); err != nil {
		return err
	}

	return nil
}