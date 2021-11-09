package repo

import (
	"context"
	"github.com/ozonmp/est-water-api/internal/model"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (
	Repo,
	context.Context,
	sqlmock.Sqlmock,
	model.Water,
) {
	mockDB, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(mockDB,"sqlmock")
	ctx := context.Background()

	repo := NewRepo(sqlxDB, uint(0))

	ts := time.Now().UTC()
	water := model.NewWater(
		uint64(0),
		"Water name",
		"Water model",
		"Water manufacturer",
		"Water material",
		uint32(100),
		&ts,
	)

	return repo, ctx, mock, *water
}

func TestDescribeWater(t *testing.T) {
	r, ctx, dbMock, dummyWater := setup(t)

	dummyWater.Id = 100
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(dummyWater.Id)

	dbMock.ExpectQuery("SELECT \\* FROM water WHERE id = \\$1").
		WithArgs(dummyWater.Id).
		WillReturnRows(rows)

	water, err := r.DescribeWater(ctx, dummyWater.Id)

	assert.NotNil(t, water)
	assert.NoError(t, err)
}

func TestListWater(t *testing.T) {
	r, ctx, dbMock, dummyWater := setup(t)

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(dummyWater.Id+1, dummyWater.Name).
		AddRow(dummyWater.Id+2, dummyWater.Name)

	dbMock.ExpectQuery("SELECT \\* FROM water ORDER BY id LIMIT 10 OFFSET 0").
		WillReturnRows(rows)

	waters, err := r.ListWaters(ctx, 10, 0)
	assert.Equal(t, 2, len(waters))
	assert.NoError(t, err)
}

func TestRemoveWater(t *testing.T) {
	r, ctx, dbMock, dummyWater := setup(t)

	dummyWater.Id = 100
	dbMock.ExpectExec("DELETE FROM water WHERE id = \\$1").
		WithArgs(dummyWater.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := r.RemoveWater(ctx, dummyWater.Id)

	assert.NoError(t, err)
}

func TestCreateWaterSuccess(t *testing.T) {
	r, ctx, dbMock, dummyWater := setup(t)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(dummyWater.Id+1)

	dbMock.ExpectQuery("INSERT INTO water \\(name,model,manufacturer,material,speed,created_at\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5,\\$6\\) RETURNING id").
		WithArgs(dummyWater.Name, dummyWater.Model, dummyWater.Manufacturer, dummyWater.Material, dummyWater.Speed, dummyWater.CreatedAt).
		WillReturnRows(rows)

	err := r.CreateWater(ctx, &dummyWater)

	assert.Equal(t, uint64(1), dummyWater.Id)
	assert.NoError(t, err)
}

func TestCreateWaterError(t *testing.T) {
	r, ctx, dbMock, dummyWater := setup(t)

	dbMock.ExpectQuery("INSERT INTO water \\(name,model,manufacturer,material,speed,created_at\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5,\\$6\\) RETURNING id").
		WithArgs(dummyWater.Name, dummyWater.Model, dummyWater.Manufacturer, dummyWater.Material, dummyWater.Speed, dummyWater.CreatedAt)

	err := r.CreateWater(ctx, &dummyWater)

	assert.Error(t, err)
}
