package producer

import (
	"context"
	"sync"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/pkg/errors"

	"github.com/ozonmp/est-water-api/internal/app/metric"
	"github.com/ozonmp/est-water-api/internal/app/repo"
	"github.com/ozonmp/est-water-api/internal/app/sender"
	"github.com/ozonmp/est-water-api/internal/logger"
	"github.com/ozonmp/est-water-api/internal/model"
)

type Producer interface {
	Start(ctx context.Context)
	Close()
}

type producer struct {
	n uint64
	timeout time.Duration

	sender sender.EventSender
	events <-chan model.WaterEvent

	workerPool *workerpool.WorkerPool
	workerBatchSize uint64
	workerBatchTimeout time.Duration

	repo repo.EventRepo

	wg *sync.WaitGroup
}

type Config struct {
	N uint64
	Sender sender.EventSender
	Events <-chan model.WaterEvent
	WorkerPool *workerpool.WorkerPool
	WorkerBatchSize uint64
	WorkerBatchTimeout time.Duration
	Repo repo.EventRepo
}

func NewKafkaProducer(cfg Config) Producer {

	wg := &sync.WaitGroup{}

	return &producer{
		n: cfg.N,
		sender: cfg.Sender,
		events: cfg.Events,
		workerPool: cfg.WorkerPool,
		workerBatchSize: cfg.WorkerBatchSize,
		workerBatchTimeout: cfg.WorkerBatchTimeout,
		repo: cfg.Repo,
		wg: wg,
	}
}

func (p *producer) Start(ctx context.Context) {
	for i := uint64(0); i < p.n; i++ {
		p.wg.Add(1)

		go func() {
			defer p.wg.Done()
			ticker := time.NewTicker(p.workerBatchTimeout)

			workerBatchUpdate := make([]uint64, 0, p.workerBatchSize)
			workerBatchClean := make([]uint64, 0, p.workerBatchSize)

			for {
				select {
				case <-ticker.C:
					p.workerBatchSendUpdate(ctx, &workerBatchUpdate)
					p.workerBatchSendClean(ctx, &workerBatchClean)
				case event := <-p.events:
					if err := p.sender.Send(&event); err != nil {
						logger.ErrorKV(ctx, "producer send event failed",
							"err", errors.Wrapf(err, "sender.Send() failed with %v", event),
						)

						workerBatchUpdate = append(workerBatchUpdate, event.ID)
						if len(workerBatchUpdate) >= int(p.workerBatchSize) {
							ticker.Reset(p.workerBatchTimeout)
							p.workerBatchSendUpdate(ctx, &workerBatchUpdate)
						}
					} else {
						workerBatchClean = append(workerBatchClean, event.ID)
						if len(workerBatchClean) >= int(p.workerBatchSize) {
							ticker.Reset(p.workerBatchTimeout)
							p.workerBatchSendClean(ctx, &workerBatchClean)
						}
					}
				case <-ctx.Done():
					ticker.Stop()
					p.workerBatchSendUpdate(ctx, &workerBatchUpdate)
					p.workerBatchSendClean(ctx, &workerBatchClean)
					p.clearUnsentEvents(ctx)
					return
				}
			}
		}()
	}
}

func (p *producer) Close() {
	p.wg.Wait()
}

func (p *producer) workerBatchSendUpdate(ctx context.Context, eventIDs *[]uint64) {
	if len(*eventIDs) > 0 {
		var ids []uint64
		for _, id := range *eventIDs {
			ids = append(ids, id)
		}

		metric.SubTotalWaterEventsNow(uint(len(ids)))

		p.workerPool.Submit(func() {
			if err := p.repo.Unlock(ctx, ids); err != nil {
				logger.ErrorKV(ctx, "producer update failed",
					"err", errors.Wrapf(err, "repo.Unlock() failed with ids=%v", ids),
				)
			}
		})
		*eventIDs = nil
	}
}

func (p *producer) workerBatchSendClean(ctx context.Context, eventIDs *[]uint64) {
	if len(*eventIDs) > 0 {
		var ids []uint64
		for _, id := range *eventIDs {
			ids = append(ids, id)
		}

		metric.SubTotalWaterEventsNow(uint(len(ids)))

		p.workerPool.Submit(func() {
			if err := p.repo.Remove(ctx, ids); err != nil {
				logger.ErrorKV(ctx, "producer clean failed",
					"err", errors.Wrapf(err, "repo.Remove() failed with ids=%v", ids),
				)
			}
		})
		*eventIDs = nil
	}
}

func (p *producer) clearUnsentEvents(ctx context.Context) {
	eventsLength := len(p.events)

	if eventsLength > 0 {
		eventIDs := make([]uint64, 0, p.workerBatchSize)
		for event := range p.events {
			eventIDs = append(eventIDs, event.ID)

			if len(eventIDs) == int(p.workerBatchSize) {
				p.workerBatchSendClean(ctx, &eventIDs)
			}
		}
		p.workerBatchSendClean(ctx, &eventIDs)
	}
}
