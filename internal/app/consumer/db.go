package consumer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/ozonmp/est-water-api/internal/app/metric"
	"github.com/ozonmp/est-water-api/internal/app/repo"
	"github.com/ozonmp/est-water-api/internal/logger"
	"github.com/ozonmp/est-water-api/internal/model"
)

type Consumer interface {
	Start(ctx context.Context)
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

	return &consumer{
		n: cfg.N,
		batchSize: cfg.BatchSize,
		timeout: cfg.Timeout,
		repo: cfg.Repo,
		events: cfg.Events,
		wg: wg,
	}
}

func (c *consumer) Start(ctx context.Context) {
	for i := uint64(0); i < c.n; i++ {
		c.wg.Add(1)

		go func() {
			defer c.wg.Done()
			ticker := time.NewTicker(c.timeout)
			for {
				select {
				case <-ticker.C:
					events, err := c.repo.Lock(ctx, c.batchSize)
					if err != nil {
						logger.ErrorKV(ctx, "consumer lock events failed",
							"err", errors.Wrap(err, "repo.Lock() failed"),
						)
						continue
					}

					var eventIDs []uint64
					for _, event := range events {
						c.events <- event
						eventIDs = append(eventIDs, event.ID)
					}

					if len(eventIDs) > 0 {
						logger.DebugKV(ctx, fmt.Sprintf("Locked eventIDs: %v", eventIDs))
					}

					totalEvents := uint(len(eventIDs))
					metric.AddTotalWaterEventsNow(totalEvents)
					metric.AddTotalWaterEvents(totalEvents)
				case <-ctx.Done():
					ticker.Stop()
					return
				}
			}
		}()
	}
}

func (c *consumer) Close() {
	c.wg.Wait()
}
