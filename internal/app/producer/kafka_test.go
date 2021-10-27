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

func TestProducer(t *testing.T) {

	testCases := map[string]struct {
		eventsCount int
		producerCount uint64
	}{
		"Success: small events count (all will be handled)": {eventsCount: 10, producerCount: 10},
		"Success: big events count (not all will be handled)": {eventsCount: 100000, producerCount: 0},
	}

	for testName, testCase := range testCases {
		testName := testName
		t.Run(testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			sender := mocks.NewMockEventSender(ctrl)
			repo := mocks.NewMockEventRepo(ctrl)

			ctx, cancel := context.WithCancel(context.Background())

			workerPool := workerpool.New(1)

			eventsCh := make(chan model.WaterEvent, testCase.eventsCount)

			for i := 0; i < testCase.eventsCount; i++ {
				eventsCh<- model.WaterEvent{}
			}

			cfg := Config{
				N: testCase.producerCount,
				Sender: sender,
				Events: eventsCh,
				WorkerPool: workerPool,
				WorkerBatchSize: 5,
				WorkerBatchTimeout: time.Millisecond,
				Repo: repo,
			}

			sender.EXPECT().Send(gomock.Eq(&model.WaterEvent{})).DoAndReturn(func(event *model.WaterEvent) error {
				return nil
			}).AnyTimes()

			repo.EXPECT().Remove(gomock.AssignableToTypeOf([]uint64{})).DoAndReturn(func(eventIDs []uint64) error {
				return nil
			}).AnyTimes()

			kafka := NewKafkaProducer(cfg)
			kafka.Start(ctx)
			time.Sleep(10*time.Millisecond)
			cancel()
			kafka.Close()

			// Проверяем, что все события обработаны
			assert.Equal(t, 0, len(eventsCh))
		})
	}

	t.Run("Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		sender := mocks.NewMockEventSender(ctrl)
		repo := mocks.NewMockEventRepo(ctrl)

		ctx, cancel := context.WithCancel(context.Background())

		workerPool := workerpool.New(1)

		eventsCount := 10
		eventsCh := make(chan model.WaterEvent, eventsCount)

		for i := 0; i < eventsCount; i++ {
			eventsCh<- model.WaterEvent{}
		}

		cfg := Config{
			N: uint64(1),
			Sender: sender,
			Events: eventsCh,
			WorkerPool: workerPool,
			WorkerBatchSize: 5,
			WorkerBatchTimeout: time.Second,
			Repo: repo,
		}

		sender.EXPECT().Send(gomock.Eq(&model.WaterEvent{})).DoAndReturn(func(event *model.WaterEvent) error {
			return errors.New("some error")
		}).AnyTimes()

		repo.EXPECT().Unlock(gomock.AssignableToTypeOf([]uint64{})).DoAndReturn(func(eventIDs []uint64) error {
			return nil
		}).AnyTimes()

		kafka := NewKafkaProducer(cfg)
		kafka.Start(ctx)
		time.Sleep(time.Millisecond)
		cancel()
		kafka.Close()

		// Проверяем, что все события обработаны
		assert.Equal(t, 0, len(eventsCh))
	})
}