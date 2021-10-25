package consumer

import (
	"sync"
	"time"

	"github.com/ozonmp/est-water-api/internal/app/repo"
	"github.com/ozonmp/est-water-api/internal/model"
)

type Consumer interface {
	Start()
	Close()
}

type consumer struct {
	n uint64
	events chan<- model.WaterEvent

	repo repo.EventRepo

	batchSize uint64
	timeout time.Duration

	done chan bool
	wg *sync.WaitGroup
}

type Config struct {
	N uint64
	Events chan<- model.WaterEvent
	Repo repo.EventRepo
	BatchSize uint64
	Timeout time.Duration
}

func NewDbConsumer(cfg Config) Consumer {

	wg := &sync.WaitGroup{}
	done := make(chan bool)

	return &consumer{
		n: cfg.N,
		batchSize: cfg.BatchSize,
		timeout: cfg.Timeout,
		repo: cfg.Repo,
		events: cfg.Events,
		wg: wg,
		done: done,
	}
}

func (c *consumer) Start() {
	for i := uint64(0); i < c.n; i++ {
		c.wg.Add(1)

		go func() {
			defer c.wg.Done()
			ticker := time.NewTicker(c.timeout)
			for {
				select {
				case <-ticker.C:
					events, err := c.repo.Lock(c.batchSize)
					if err != nil {
						continue
					}
					for _, event := range events {
						if event.Type == model.Created {
							c.events <- event
						}
					}
				case <-c.done:
					return
				}
			}
		}()
	}
}

func (c *consumer) Close() {
	close(c.done)
	c.wg.Wait()
}
