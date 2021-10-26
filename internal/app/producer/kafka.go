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
	Start()
	Close()
}

type producer struct {
	ctx context.Context

	n uint64
	timeout time.Duration

	sender sender.EventSender
	events <-chan model.WaterEvent

	workerPool *workerpool.WorkerPool
	workerBatchSize uint64
	workerBatchTimeout time.Duration
	workerBatchUpdateCh chan uint64
	workerBatchCleanCh chan uint64

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

func NewKafkaProducer(ctx context.Context, cfg Config) Producer {

	wg := &sync.WaitGroup{}

	workerBatchUpdateCh := make(chan uint64, cfg.WorkerBatchSize)
	workerBatchCleanCh := make(chan uint64, cfg.WorkerBatchSize)

	return &producer{
		ctx: ctx,
		n: cfg.N,
		sender: cfg.Sender,
		events: cfg.Events,
		workerPool: cfg.WorkerPool,
		workerBatchSize: cfg.WorkerBatchSize,
		workerBatchTimeout: cfg.WorkerBatchTimeout,
		workerBatchUpdateCh: workerBatchUpdateCh,
		workerBatchCleanCh: workerBatchCleanCh,
		repo: cfg.Repo,
		wg: wg,
	}
}

func (p *producer) Start() {
	p.wg.Add(1)

	eventSentNotifyCh := make(chan bool, 1)
	go p.workerBatchSend(eventSentNotifyCh)

	for i := uint64(0); i < p.n; i++ {
		p.wg.Add(1)

		go func() {
			defer p.wg.Done()
			for {
				select {
				case event := <-p.events:
					if err := p.sender.Send(&event); err != nil {
						log.Printf("EventSender Send event error: %v\n", err)

						p.workerBatchUpdateCh <- event.ID
						eventSentNotifyCh<- true
					} else {
						p.workerBatchCleanCh <- event.ID
						eventSentNotifyCh<- true
					}
				case <-p.ctx.Done():
					return
				}
			}
		}()
	}
}

func (p *producer) Close() {
	p.wg.Wait()
	close(p.workerBatchUpdateCh)
	close(p.workerBatchCleanCh)
}

func (p *producer) workerBatchSendUpdate() {
	batchLength := len(p.workerBatchUpdateCh)

	if batchLength > 0 {
		eventIDs := make([]uint64, 0, batchLength)
		for k := 0; k < batchLength; k++ {
			eventIDs = append(eventIDs, <-p.workerBatchUpdateCh)
		}

		p.workerPool.Submit(func() {
			if err := p.repo.Unlock(eventIDs); err != nil {
				log.Printf("EventRepo Unlock events error: %v\n", err)
			}
		})
	}
}

func (p *producer) workerBatchSendClean() {
	batchLength := len(p.workerBatchCleanCh)

	if batchLength > 0 {
		eventIDs := make([]uint64, 0, batchLength)
		for k := 0; k < batchLength; k++ {
			eventIDs = append(eventIDs, <-p.workerBatchCleanCh)
		}

		p.workerPool.Submit(func() {
			if err := p.repo.Remove(eventIDs); err != nil {
				log.Printf("EventRepo Remove events error: %v\n", err)
			}
		})
	}
}

func (p *producer) workerBatchSend(eventSentNotifyCh chan bool)  {
	defer p.wg.Done()
	ticker := time.NewTicker(p.workerBatchTimeout)
	for {
		select {
		case <-eventSentNotifyCh:
			if len(p.workerBatchUpdateCh) == int(p.workerBatchSize) {
				ticker.Reset(p.workerBatchTimeout)
				p.workerBatchSendUpdate()
			}

			if len(p.workerBatchCleanCh) == int(p.workerBatchSize) {
				ticker.Reset(p.workerBatchTimeout)
				p.workerBatchSendClean()
			}
		// Отправим события, если долго не приходят новые
		case <-ticker.C:
			p.workerBatchSendUpdate()
			p.workerBatchSendClean()
		// Завершаем необработанные события
		case <-p.ctx.Done():
			p.workerBatchSendUpdate()
			p.workerBatchSendClean()
			return
		}
	}
}
