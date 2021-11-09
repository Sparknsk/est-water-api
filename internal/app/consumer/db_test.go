package consumer

import (
	"context"
	"testing"
	"time"

	"github.com/ozonmp/est-water-api/internal/mocks"
	"github.com/ozonmp/est-water-api/internal/model"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestConsumerSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockEventRepo(ctrl)

	ctx, cancel := context.WithCancel(context.Background())

	batchSize := uint64(1)
	consumerCount := 5
	eventsCh := make(chan model.WaterEvent, consumerCount-2)

	ts := time.Now().UTC()
	dummyEvent := model.WaterEvent{
		ID: uint64(123),
		WaterId: uint64(123),
		Type: model.Created,
		Status: model.Unlocked,
		Entity: model.NewWater(
			uint64(123),
			"name",
			"model",
			"manufacturer",
			"material",
			uint32(100),
			&ts,
			false,
		),
	}

	repo.EXPECT().Lock(ctx, gomock.Eq(batchSize)).DoAndReturn(func(ctx context.Context, n uint64) ([]model.WaterEvent, error) {
		return []model.WaterEvent{dummyEvent}, nil
	}).Times(consumerCount)

	cfg := Config{
		N: uint64(consumerCount),
		Events: eventsCh,
		Repo: repo,
		BatchSize: batchSize,
		Timeout: time.Millisecond*100,
	}

	db := NewDbConsumer(cfg)

	db.Start(ctx)

	// Проверяем, что в канале находится исходное событие
	for i := 0; i < consumerCount; i++ {
		assert.Equal(t, dummyEvent, <-eventsCh)
	}
	// Проверяем, что количество событий в канале было равно количеству потребителей
	assert.Equal(t, 0, len(eventsCh))

	cancel()

	db.Close()
}