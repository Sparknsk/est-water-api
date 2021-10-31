package retranslator

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ozonmp/est-water-api/internal/mocks"
	"github.com/ozonmp/est-water-api/internal/model"

	"github.com/golang/mock/gomock"
)

func setup(t *testing.T, eventsCount int) (
	*mocks.MockEventRepo,
	*mocks.MockEventSender,
	context.Context,
	context.CancelFunc,
	[]model.WaterEvent,
) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)
	ctx, cancel := context.WithCancel(context.Background())

	dummyEvent := model.WaterEvent{
		ID: uint64(1),
		Type: model.Created,
		Status: model.Processed,
		Entity: model.NewWater(
			uint64(1),
			"name",
			"model",
			"manufacturer",
			"material",
			uint32(100),
		),
	}
	events := make([]model.WaterEvent, 0, eventsCount)
	for i := 0; i < eventsCount; i++ {
		events = append(events, dummyEvent)
	}

	return repo, sender, ctx, cancel, events
}

func TestRetranslator(t *testing.T) {

	testCases := map[string]struct {
		eventsCount int
		consumerCount uint64
		producerCount uint64
	}{
		"Success: small events count": {eventsCount: 20, consumerCount: 2, producerCount: 2},
		"Success: big events count (consumers > producers)": {eventsCount: 2007, consumerCount: 15, producerCount: 5},
		"Success: big events count (consumers < producers)": {eventsCount: 2007, consumerCount: 5, producerCount: 15},
		"Success: huge events count (consumers < producers)": {eventsCount: 20234, consumerCount: 50, producerCount: 150},
		"Success: huge events count (consumers > producers)": {eventsCount: 20234, consumerCount: 150, producerCount: 50},
	}

	for testName, testCase := range testCases {
		testName := testName
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			repo, sender, ctx, cancel, events := setup(t, testCase.eventsCount)

			cfg := Config {
				ChannelSize: 512,
				ConsumerCount: testCase.consumerCount,
				ConsumeSize: 10,
				ConsumeTimeout: time.Millisecond*2,
				ProducerCount: testCase.producerCount,
				WorkerCount: 1,
				WorkerBatchSize: 5,
				WorkerBatchTimeout: time.Second,
				Repo: repo,
				Sender: sender,
			}

			eventsDone := make(chan bool, 1)
			startId := uint64(0)
			stopId := cfg.ConsumeSize
			repo.EXPECT().Lock(gomock.Eq(cfg.ConsumeSize)).DoAndReturn(func(n uint64) ([]model.WaterEvent, error) {
				time.Sleep(time.Millisecond*100)

				start := atomic.LoadUint64(&startId)
				atomic.AddUint64(&startId, n)
				stop := atomic.LoadUint64(&stopId)
				atomic.AddUint64(&stopId, n)

				if stop > uint64(testCase.eventsCount) {
					stop = uint64(testCase.eventsCount)
				}

				if start > stop || testCase.eventsCount == 0 {
					if len(eventsDone) < 1 {
						eventsDone <- true
					}
					return nil, nil
				}

				return events[start:stop], nil
			}).AnyTimes()

			sender.EXPECT().Send(&events[0]).DoAndReturn(func(event *model.WaterEvent) error {
				return nil
			}).AnyTimes()

			repo.EXPECT().Remove(gomock.AssignableToTypeOf([]uint64{})).DoAndReturn(func(eventIDs []uint64) error {
				return nil
			}).AnyTimes()

			repo.EXPECT().Unlock(gomock.AssignableToTypeOf([]uint64{})).DoAndReturn(func(eventIDs []uint64) error {
				return nil
			}).AnyTimes()

			transponder := NewRetranslator(cfg)
			transponder.Start(ctx)
			<-eventsDone
			cancel()
			transponder.Close()
		})
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		repo, sender, ctx, cancel, events := setup(t, 10)

		cfg := Config {
			ChannelSize: 512,
			ConsumerCount: 1,
			ConsumeSize: 10,
			ConsumeTimeout: time.Millisecond*2,
			ProducerCount: 1,
			WorkerCount: 1,
			WorkerBatchSize: 5,
			WorkerBatchTimeout: time.Second,
			Repo: repo,
			Sender: sender,
		}

		repo.EXPECT().Lock(gomock.Eq(cfg.ConsumeSize)).DoAndReturn(func(n uint64) ([]model.WaterEvent, error) {
			time.Sleep(time.Millisecond*100)
			return events, nil
		}).AnyTimes()

		sender.EXPECT().Send(gomock.Any()).DoAndReturn(func(event *model.WaterEvent) error {
			return errors.New("some error")
		}).AnyTimes()

		repo.EXPECT().Unlock(gomock.AssignableToTypeOf([]uint64{})).DoAndReturn(func(eventIDs []uint64) error {
			return nil
		}).AnyTimes()

		transponder := NewRetranslator(cfg)
		transponder.Start(ctx)
		time.Sleep(time.Millisecond)
		cancel()
		transponder.Close()
	})
}
