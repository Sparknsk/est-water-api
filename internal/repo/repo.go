package repo

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/ozonmp/est-water-api/internal/database"
	"github.com/ozonmp/est-water-api/internal/model"
)

const waterTableName = "water"

type Repo interface {
	DescribeWater(ctx context.Context, waterID uint64) (*model.Water, error)
	CreateWater(ctx context.Context, water *model.Water) error
	ListWaters(ctx context.Context, limit uint64, offset uint64) ([]model.Water, error)
	RemoveWater(ctx context.Context, waterID uint64) error
}

type repo struct {
	db *sqlx.DB
	batchSize uint
}

func NewRepo(db *sqlx.DB, batchSize uint) Repo {
	return &repo{db: db, batchSize: batchSize}
}

func (r *repo) DescribeWater(ctx context.Context, waterID uint64) (*model.Water, error) {
	query := database.StatementBuilder.
		Select("*").
		From(waterTableName).
		Where(sq.Eq{"id": waterID, "delete_status": false})

	queryText, queryArgs, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryxContext(ctx, queryText, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var water model.Water
	for rows.Next() {
		if err = rows.StructScan(&water); err != nil {
			return nil, err
		}
	}

	if water.Id == 0 {
		return nil, nil
	}

	return &water, nil
}

func (r *repo) CreateWater(ctx context.Context, water *model.Water) error {
	query := database.StatementBuilder.
		Insert(waterTableName).
		Columns("name", "model", "manufacturer", "material", "speed", "created_at").
		Values(water.Name, water.Model, water.Manufacturer, water.Material, water.Speed, water.CreatedAt).
		Suffix("RETURNING id").
		RunWith(r.db)

	if err := query.QueryRowContext(ctx).Scan(&water.Id); err != nil {
		return err
	}

	return nil
}

func (r *repo) ListWaters(ctx context.Context, limit uint64, offset uint64) ([]model.Water, error) {
	query := database.StatementBuilder.
		Select("*").
		From(waterTableName).
		Where(sq.Eq{"delete_status": false}).
		Limit(limit).
		Offset(offset).
		OrderBy("id")

	queryText, queryArgs, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var res []model.Water
	rows, err := r.db.QueryxContext(ctx, queryText, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var water model.Water
	for rows.Next() {
		if err = rows.StructScan(&water); err != nil {
			return nil, err
		}
		res = append(res, water)
	}

	return res, err
}

func (r *repo) RemoveWater(ctx context.Context, waterID uint64) error {
	query := database.StatementBuilder.
		Update(waterTableName).
		Set("delete_status", true).
		Where(sq.Eq{"id": waterID})

	queryText, queryArgs, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, queryText, queryArgs...)
	return err
}
