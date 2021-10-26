package retranslator

import (
	"context"
	"errors"
	"math"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ozonmp/est-water-api/internal/mocks"
	"github.com/ozonmp/est-water-api/internal/model"

	"github.com/golang/mock/gomock"
)

var dummyEvent = model.WaterEvent{
	ID: uint64(123),
	Type: model.Created,
	Status: model.Processed,
	Entity: model.NewWater(
		uint64(123),
		"name",
		"model",
		"manufacturer",
		"material",
		100,
	),
}

func TestRetranslatorSuccess(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repo := mocks.NewMockEventRepo(ctrl)

	sender := mocks.NewMockEventSender(ctrl)

	ctx, cancel := context.WithCancel(context.Background())

	eventsCount := 28
	var allEvents []model.WaterEvent
	for i := 1; i <= eventsCount; i++ {
		allEvents = append(allEvents, dummyEvent)
	}

	cfg := Config {
		ChannelSize: 512,
		ConsumerCount: 2,
		ConsumeSize: 10,
		ConsumeTimeout: time.Second,
		ProducerCount: 3,
		WorkerCount: 1,
		WorkerBatchSize: 5,
		WorkerBatchTimeout: time.Second*10,
		Repo: repo,
		Sender: sender,
	}

	startId := int32(0)
	stopId := int32(cfg.ConsumeSize)
	repo.EXPECT().Lock(gomock.Eq(cfg.ConsumeSize)).DoAndReturn(func(n uint64) ([]model.WaterEvent, error) {
		start := atomic.LoadInt32(&startId)
		atomic.AddInt32(&startId, int32(n))
		stop := atomic.LoadInt32(&stopId)
		atomic.AddInt32(&stopId, int32(n))

		if stop > int32(eventsCount) {
			stop = int32(eventsCount)
		}

		if start > stop {
			return nil, nil
		}

		return allEvents[start:stop], nil
	}).Times(int(2*cfg.ConsumerCount))

	sender.EXPECT().Send(gomock.Eq(&dummyEvent)).DoAndReturn(func(event *model.WaterEvent) error {
		return nil
	}).Times(eventsCount)

	// Remove будет вызван в зависимости от количества событий eventsCount полученных из консьюмера
	// разбитых на пачки по cfg.WorkerBatchSize
	removeTimes := int(math.Round(float64(eventsCount)/float64(cfg.WorkerBatchSize)))
	repo.EXPECT().Remove(gomock.AssignableToTypeOf([]uint64{})).DoAndReturn(func(eventIDs []uint64) error {
		return nil
	}).Times(removeTimes)

	transponder := NewRetranslator(cfg)
	transponder.Start(ctx)
	time.Sleep(2*cfg.ConsumeTimeout + 100*time.Millisecond)
	cancel()
	transponder.Close()
}

func TestRetranslatorError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repo := mocks.NewMockEventRepo(ctrl)

	sender := mocks.NewMockEventSender(ctrl)

	ctx, cancel := context.WithCancel(context.Background())

	cfg := Config {
		ChannelSize: 512,
		ConsumerCount: 2,
		ConsumeSize: 10,
		ConsumeTimeout: time.Second,
		ProducerCount: 3,
		WorkerCount: 1,
		WorkerBatchSize: 5,
		WorkerBatchTimeout: time.Millisecond*100,
		Repo: repo,
		Sender: sender,
	}

	repo.EXPECT().Lock(gomock.Eq(cfg.ConsumeSize)).DoAndReturn(func(n uint64) ([]model.WaterEvent, error) {
		return []model.WaterEvent{dummyEvent}, nil
	}).Times(int(cfg.ConsumerCount))

	sender.EXPECT().Send(gomock.Eq(&dummyEvent)).DoAndReturn(func(event *model.WaterEvent) error {
		return errors.New("some error")
	}).Times(int(cfg.ConsumerCount))

	repo.EXPECT().Unlock(gomock.AssignableToTypeOf([]uint64{})).DoAndReturn(func(eventIDs []uint64) error {
		return nil
	}).AnyTimes()

	transponder := NewRetranslator(cfg)
	transponder.Start(ctx)
	time.Sleep(cfg.ConsumeTimeout + 100*time.Millisecond)
	cancel()
	transponder.Close()
}
