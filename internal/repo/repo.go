package repo

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/ozonmp/est-water-api/internal/database"
	"github.com/ozonmp/est-water-api/internal/model"
)

const waterTableName = "water"

//go:generate mockgen -destination=../mocks/repo_mock.go -package=mocks github.com/ozonmp/est-water-api/internal/repo Repo
type Repo interface {
	Get(ctx context.Context, waterId uint64) (*model.Water, error)
	Create(ctx context.Context, water *model.Water) error
	List(ctx context.Context, limit uint64, offset uint64) ([]model.Water, error)
	Remove(ctx context.Context, waterId uint64) error
	Update(ctx context.Context, water *model.Water) error
}

type repo struct {
	db *sqlx.DB
	batchSize uint
}

func NewRepo(db *sqlx.DB, batchSize uint) Repo {
	return &repo{db: db, batchSize: batchSize}
}

func (r *repo) Get(ctx context.Context, waterId uint64) (*model.Water, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.Get()")
	defer span.Finish()
	span.LogKV(
		"event", "water repository get",
		"waterId", waterId,
	)

	query := database.StatementBuilder.
		Select("*").
		From(waterTableName).
		Where(sq.Eq{"id": waterId, "delete_status": false})

	queryText, queryArgs, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "query.ToSql() failed")
	}

	rows, err := r.db.QueryxContext(ctx, queryText, queryArgs...)
	if err != nil {
		return nil, errors.Wrap(err, "db.QueryxContext() failed")
	}
	defer rows.Close()

	var water model.Water
	for rows.Next() {
		if err = rows.StructScan(&water); err != nil {
			return nil, errors.Wrap(err, "rows.StructScan() failed")
		}
	}

	if water.Id == 0 {
		return nil, nil
	}

	return &water, nil
}

func (r *repo) Create(ctx context.Context, water *model.Water) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.Create()")
	defer span.Finish()
	span.LogKV(
		"event", "water repository create",
		"water", water,
	)

	query := database.StatementBuilder.
		Insert(waterTableName).
		Columns("name", "model", "manufacturer", "material", "speed", "created_at").
		Values(water.Name, water.Model, water.Manufacturer, water.Material, water.Speed, water.CreatedAt).
		Suffix("RETURNING id").
		RunWith(r.db)

	if err := query.QueryRowContext(ctx).Scan(&water.Id); err != nil {
		return errors.Wrap(err, "query.QueryRowContext().Scan() failed")
	}

	return nil
}

func (r *repo) List(ctx context.Context, limit uint64, offset uint64) ([]model.Water, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.List()")
	defer span.Finish()
	span.LogKV(
		"event", "water repository list",
		"limit", limit,
		"offset", offset,
	)

	query := database.StatementBuilder.
		Select("*").
		From(waterTableName).
		Where(sq.Eq{"delete_status": false}).
		Limit(limit).
		Offset(offset).
		OrderBy("id")

	queryText, queryArgs, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "query.ToSql() failed")
	}

	var res []model.Water
	rows, err := r.db.QueryxContext(ctx, queryText, queryArgs...)
	if err != nil {
		return nil, errors.Wrap(err, "db.QueryxContext() failed")
	}
	defer rows.Close()

	var water model.Water
	for rows.Next() {
		if err = rows.StructScan(&water); err != nil {
			return nil, errors.Wrap(err, "rows.StructScan() failed")
		}
		res = append(res, water)
	}

	return res, err
}

func (r *repo) Remove(ctx context.Context, waterId uint64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.Remove()")
	defer span.Finish()
	span.LogKV(
		"event", "water repository remove",
		"waterId", waterId,
	)

	query := database.StatementBuilder.
		Update(waterTableName).
		Set("delete_status", true).
		Where(sq.Eq{"id": waterId}).
		RunWith(r.db)

	_, err := query.ExecContext(ctx)
	return errors.Wrap(err, "query.ExecContext() failed")
}

func (r *repo) Update(ctx context.Context, water *model.Water) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.Update()")
	defer span.Finish()
	span.LogKV(
		"event", "water repository update",
		"water", water,
	)

	query := database.StatementBuilder.
		Update(waterTableName).
		Set("name", water.Name).
		Set("model", water.Model).
		Set("material", water.Material).
		Set("manufacturer", water.Manufacturer).
		Set("speed", water.Speed).
		Set("updated_at", water.UpdatedAt).
		Where(sq.Eq{"id": water.Id}).
		RunWith(r.db)

	_, err := query.ExecContext(ctx)
	return errors.Wrap(err, "query.ExecContext() failed")
}
