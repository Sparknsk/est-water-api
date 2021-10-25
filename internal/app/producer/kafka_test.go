package producer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ozonmp/est-water-api/internal/mocks"
	"github.com/ozonmp/est-water-api/internal/model"

	"github.com/gammazero/workerpool"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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

func TestProducerSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockEventSender(ctrl)
	repo := mocks.NewMockEventRepo(ctrl)

	eventsCount := 7
	eventsCh := make(chan model.WaterEvent, eventsCount-2)

	workerPool := workerpool.New(1)

	sender.EXPECT().Send(gomock.Eq(&dummyEvent)).DoAndReturn(func(event *model.WaterEvent) error {
		return nil
	}).Times(eventsCount)

	repo.EXPECT().Remove(gomock.Eq([]uint64{dummyEvent.ID})).DoAndReturn(func(eventIDs []uint64) error {
		return nil
	}).Times(eventsCount)

	cfg := Config{
		N: uint64(3),
		Sender: sender,
		Events: eventsCh,
		WorkerPool: workerPool,
		Repo: repo,
	}

	ctx, cancel := context.WithCancel(context.Background())

	kafka := NewKafkaProducer(ctx, cfg)

	kafka.Start()

	for i := 0; i < eventsCount; i++ {
		eventsCh<- dummyEvent
	}

	time.Sleep(100*time.Millisecond)

	cancel()

	kafka.Close()

	// Проверяем, что все события обработаны
	assert.Equal(t, 0, len(eventsCh))
}

func TestProducerError(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockEventSender(ctrl)
	repo := mocks.NewMockEventRepo(ctrl)

	eventsCount := 7
	eventsCh := make(chan model.WaterEvent, eventsCount-2)

	workerPool := workerpool.New(1)

	sender.EXPECT().Send(gomock.Eq(&dummyEvent)).DoAndReturn(func(event *model.WaterEvent) error {
		return errors.New("some error")
	}).Times(eventsCount)

	repo.EXPECT().Unlock(gomock.Eq([]uint64{dummyEvent.ID})).DoAndReturn(func(eventIDs []uint64) error {
		return nil
	}).Times(eventsCount)

	cfg := Config{
		N: uint64(3),
		Sender: sender,
		Events: eventsCh,
		WorkerPool: workerPool,
		Repo: repo,
	}

	ctx, cancel := context.WithCancel(context.Background())

	kafka := NewKafkaProducer(ctx, cfg)

	kafka.Start()

	for i := 0; i < eventsCount; i++ {
		eventsCh<- dummyEvent
	}

	time.Sleep(100*time.Millisecond)

	cancel()

	kafka.Close()

	// Проверяем, что все события обработаны
	assert.Equal(t, 0, len(eventsCh))
}