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
	EventRepo,
	context.Context,
	sqlmock.Sqlmock,
	model.WaterEvent,
) {
	mockDB, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(mockDB,"sqlmock")
	ctx := context.Background()

	repo := NewEventRepo(sqlxDB)

	ts := time.Now().UTC()
	waterEvent := model.WaterEvent{
		ID: uint64(0),
		WaterId: uint64(1),
		Type: model.Created,
		Status: model.Unlocked,
		Entity: model.NewWater(
			uint64(1),
			"name",
			"model",
			"manufacturer",
			"material",
			uint32(100),
			&ts,
			false,
		),
	}

	return repo, ctx, mock, waterEvent
}

func TestLock(t *testing.T) {
	r, ctx, dbMock, dummyWaterEvent := setup(t)

	rows := sqlmock.NewRows([]string{"id", "water_id", "type", "status", "payload"}).
		AddRow(dummyWaterEvent.ID+1, dummyWaterEvent.WaterId, dummyWaterEvent.Type, dummyWaterEvent.Status, *dummyWaterEvent.Entity).
		AddRow(dummyWaterEvent.ID+2, dummyWaterEvent.WaterId+1, dummyWaterEvent.Type, dummyWaterEvent.Status, *dummyWaterEvent.Entity)

	dbMock.ExpectQuery("WITH cte AS \\( SELECT id FROM water_events WHERE status = \\$1 ORDER BY id LIMIT 2 FOR NO KEY UPDATE \\) UPDATE water_events we SET status = \\$2 FROM cte WHERE we.id = cte.id RETURNING we.\\*").
		WithArgs("unlock", "lock").
		WillReturnRows(rows)

	events, err := r.Lock(ctx, 2)

	assert.NotEmpty(t, events)
	assert.NoError(t, err)
}

func TestUnlock(t *testing.T) {
	r, ctx, dbMock, _ := setup(t)

	eventIDs := []uint64{1, 2, 3}
	dbMock.ExpectExec("UPDATE water_events SET status = \\$1, updated_at = \\$2 WHERE id IN \\(\\$3,\\$4,\\$5\\)").
		WithArgs("unlock", sqlmock.AnyArg(), eventIDs[0], eventIDs[1], eventIDs[2]).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := r.Unlock(ctx, eventIDs)

	assert.NoError(t, err)
}

func TestRemove(t *testing.T) {
	r, ctx, dbMock, _ := setup(t)

	eventIDs := []uint64{1, 2, 3}
	dbMock.ExpectExec("DELETE FROM water_events WHERE id IN \\(\\$1,\\$2,\\$3\\)").
		WithArgs(eventIDs[0], eventIDs[1], eventIDs[2]).
		WillReturnResult(sqlmock.NewResult(0, 3))

	err := r.Remove(ctx, eventIDs)

	assert.NoError(t, err)
}

func TestAdd(t *testing.T) {
	r, ctx, dbMock, dummyWaterEvent := setup(t)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(dummyWaterEvent.ID+1)

	dbMock.ExpectQuery("INSERT INTO water_events \\(water_id,type,status,payload,created_at\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5\\) RETURNING id").
		WithArgs(dummyWaterEvent.WaterId, dummyWaterEvent.Type, dummyWaterEvent.Status, dummyWaterEvent.Entity, dummyWaterEvent.CreatedAt).
		WillReturnRows(rows)

	err := r.Add(ctx, &dummyWaterEvent)

	assert.Equal(t, uint64(1), dummyWaterEvent.ID)
	assert.NoError(t, err)
}