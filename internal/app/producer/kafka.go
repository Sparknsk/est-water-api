package producer

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ozonmp/est-water-api/internal/app/repo"
	"github.com/ozonmp/est-water-api/internal/app/sender"
	"github.com/ozonmp/est-water-api/internal/model"

	"github.com/gammazero/workerpool"
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
					p.workerBatchSendUpdate(&workerBatchUpdate)
					p.workerBatchSendClean(&workerBatchClean)
				case event := <-p.events:
					if err := p.sender.Send(&event); err != nil {
						log.Printf("EventSender Send event error: %v\n", err)

						workerBatchUpdate = append(workerBatchUpdate, event.ID)
						if len(workerBatchUpdate) == int(p.workerBatchSize) {
							ticker.Reset(p.workerBatchTimeout)
							p.workerBatchSendUpdate(&workerBatchUpdate)
						}
					} else {
						workerBatchClean = append(workerBatchClean, event.ID)
						if len(workerBatchClean) == int(p.workerBatchSize) {
							ticker.Reset(p.workerBatchTimeout)
							p.workerBatchSendClean(&workerBatchClean)
						}
					}
				case <-ctx.Done():
					ticker.Stop()
					p.workerBatchSendUpdate(&workerBatchUpdate)
					p.workerBatchSendClean(&workerBatchClean)
					return
				}
			}
		}()
	}
}

func (p *producer) Close() {
	p.wg.Wait()
	p.clearUnsentEvents()
}

func (p *producer) workerBatchSendUpdate(eventIDs *[]uint64) {
	if len(*eventIDs) > 0 {
		var ids []uint64
		for _, id := range *eventIDs {
			ids = append(ids, id)
		}
		p.workerPool.Submit(func() {
			if err := p.repo.Unlock(ids); err != nil {
				log.Printf("EventRepo Unlock events error: %v\n", err)
			}
		})
		*eventIDs = nil
	}
}

func (p *producer) workerBatchSendClean(eventIDs *[]uint64) {
	if len(*eventIDs) > 0 {
		var ids []uint64
		for _, id := range *eventIDs {
			ids = append(ids, id)
		}
		p.workerPool.Submit(func() {
			if err := p.repo.Remove(ids); err != nil {
				log.Printf("EventRepo Remove events error: %v\n", err)
			}
		})
		*eventIDs = nil
	}
}

func (p *producer) clearUnsentEvents() {
	eventsLength := len(p.events)

	if eventsLength > 0 {
		eventIDs := make([]uint64, 0, p.workerBatchSize)
		eventsCounter := 1
		for i := 0; i < eventsLength; i++ {
			event := <-p.events

			eventsCounter++
			eventIDs = append(eventIDs, event.ID)

			if len(eventIDs) == int(p.workerBatchSize) || eventsCounter == eventsLength {
				p.workerBatchSendClean(&eventIDs)
			}
		}
	}
}
