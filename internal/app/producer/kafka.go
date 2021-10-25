package producer

import (
	"github.com/ozonmp/est-water-api/internal/app/repo"
	"log"
	"sync"
	"time"

	"github.com/ozonmp/est-water-api/internal/app/sender"
	"github.com/ozonmp/est-water-api/internal/model"

	"github.com/gammazero/workerpool"
)

type Producer interface {
	Start()
	Close()
}

type producer struct {
	n       uint64
	timeout time.Duration

	sender sender.EventSender
	events <-chan model.WaterEvent

	workerPool *workerpool.WorkerPool

	repo repo.EventRepo

	wg   *sync.WaitGroup
	done chan bool
}

type Config struct {
	N uint64
	Sender sender.EventSender
	Events <-chan model.WaterEvent
	WorkerPool *workerpool.WorkerPool
	Repo repo.EventRepo
}

func NewKafkaProducer(cfg Config) Producer {

	wg := &sync.WaitGroup{}
	done := make(chan bool)

	return &producer{
		n: cfg.N,
		sender: cfg.Sender,
		events: cfg.Events,
		workerPool: cfg.WorkerPool,
		repo: cfg.Repo,
		wg: wg,
		done: done,
	}
}

func (p *producer) Start() {
	for i := uint64(0); i < p.n; i++ {
		p.wg.Add(1)

		go func() {
			defer p.wg.Done()
			for {
				select {
				case event := <-p.events:
					if err := p.sender.Send(&event); err != nil {
						//log.Printf("EventSender Send event error: %v\n", err)

						p.workerPool.Submit(func() {
							if err := p.repo.Unlock([]uint64{event.ID}); err != nil {
								log.Printf("EventRepo Unlock error: %v\n", err)
							}
						})
					} else {
						p.workerPool.Submit(func() {
							if err := p.repo.Remove([]uint64{event.ID}); err != nil {
								log.Printf("EventRepo Remove error: %v\n", err)
							}
						})
					}
				case <-p.done:
					return
				}
			}
		}()
	}
}

func (p *producer) Close() {
	close(p.done)
	p.wg.Wait()
}
