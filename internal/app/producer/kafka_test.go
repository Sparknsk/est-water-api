package producer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/ozonmp/est-water-api/internal/logger"
	"github.com/ozonmp/est-water-api/internal/mocks"
	"github.com/ozonmp/est-water-api/internal/model"
)

func setup(t *testing.T, eventsCount int) (
	*mocks.MockEventRepo,
	*mocks.MockEventSender,
	*workerpool.WorkerPool,
	context.Context,
	context.CancelFunc,
	chan model.WaterEvent,
) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)
	workerPool := workerpool.New(1)
	ctx, cancel := context.WithCancel(context.Background())

	newLogger := logger.CloneWithLevel(ctx, zap.FatalLevel)
	ctx = logger.AttachLogger(ctx, newLogger)

	eventsCh := make(chan model.WaterEvent, eventsCount)

	for i := 0; i < eventsCount; i++ {
		eventsCh<- model.WaterEvent{}
	}
	close(eventsCh)

	return repo, sender, workerPool, ctx, cancel, eventsCh
}

func TestProducer(t *testing.T) {

	testCases := map[string]struct {
		eventsCount int
		producerCount uint64
	}{
		"Success: small events count (all will be handled - batch length)": {eventsCount: 200, producerCount: 2},
		"Success: small events count (all will be handled - batch timeout)": {eventsCount: 50, producerCount: 30},
		"Success: big events count": {eventsCount: 10000, producerCount: 1},
	}

	for testName, testCase := range testCases {
		testName := testName
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			repo, sender, workerPool, ctx, cancel, eventsCh := setup(t, testCase.eventsCount)

			cfg := Config{
				N: testCase.producerCount,
				Sender: sender,
				Events: eventsCh,
				WorkerPool: workerPool,
				WorkerBatchSize: 50,
				WorkerBatchTimeout: time.Millisecond,
				Repo: repo,
			}

			sender.EXPECT().Send(ctx, gomock.Eq(&model.WaterEvent{})).DoAndReturn(func(ctx context.Context, event *model.WaterEvent) error {
				return nil
			}).AnyTimes()

			repo.EXPECT().Remove(ctx, gomock.AssignableToTypeOf([]uint64{})).DoAndReturn(func(ctx context.Context, eventIDs []uint64) error {
				return nil
			}).AnyTimes()

			kafka := NewKafkaProducer(cfg)
			kafka.Start(ctx)
			time.Sleep(100*time.Millisecond)
			cancel()
			kafka.Close()

			// Проверяем, что все события обработаны
			assert.Equal(t, 0, len(eventsCh))
		})
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		repo, sender, workerPool, ctx, cancel, eventsCh := setup(t, 10)

		cfg := Config{
			N: uint64(1),
			Sender: sender,
			Events: eventsCh,
			WorkerPool: workerPool,
			WorkerBatchSize: 5,
			WorkerBatchTimeout: time.Second,
			Repo: repo,
		}

		sender.EXPECT().Send(ctx, gomock.Eq(&model.WaterEvent{})).DoAndReturn(func(ctx context.Context, event *model.WaterEvent) error {
			return errors.New("some error")
		}).AnyTimes()

		repo.EXPECT().Unlock(ctx, gomock.AssignableToTypeOf([]uint64{})).DoAndReturn(func(ctx context.Context, eventIDs []uint64) error {
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