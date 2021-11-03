package retranslator

import (
	"context"
	"time"

	"github.com/ozonmp/est-water-api/internal/app/consumer"
	"github.com/ozonmp/est-water-api/internal/app/producer"
	"github.com/ozonmp/est-water-api/internal/app/repo"
	"github.com/ozonmp/est-water-api/internal/app/sender"
	"github.com/ozonmp/est-water-api/internal/model"

	"github.com/gammazero/workerpool"
)

type Retranslator interface {
	Start(ctx context.Context)
	Close()
}

type Config struct {
	ChannelSize uint64

	ConsumerCount uint64
	ConsumeSize uint64
	ConsumeTimeout time.Duration

	ProducerCount uint64
	WorkerCount int
	WorkerBatchSize uint64
	WorkerBatchTimeout time.Duration

	Repo repo.EventRepo
	Sender sender.EventSender
}

type retranslator struct {
	events chan model.WaterEvent
	consumer consumer.Consumer
	producer producer.Producer
	workerPool *workerpool.WorkerPool
}

func NewRetranslator(cfg Config) Retranslator {
	events := make(chan model.WaterEvent, cfg.ChannelSize)

	workerPool := workerpool.New(cfg.WorkerCount)

	consumerCfg := consumer.Config{
		N: cfg.ConsumerCount,
		Events: events,
		Repo: cfg.Repo,
		BatchSize: cfg.ConsumeSize,
		Timeout: cfg.ConsumeTimeout,
	}

	consumer := consumer.NewDbConsumer(consumerCfg)

	producerCfg := producer.Config{
		N: cfg.ProducerCount,
		Sender: cfg.Sender,
		Events: events,
		WorkerPool: workerPool,
		WorkerBatchSize: cfg.WorkerBatchSize,
		WorkerBatchTimeout: cfg.WorkerBatchTimeout,
		Repo: cfg.Repo,
	}

	producer := producer.NewKafkaProducer(producerCfg)

	return &retranslator{
		events: events,
		consumer: consumer,
		producer: producer,
		workerPool: workerPool,
	}
}

func (r *retranslator) Start(ctx context.Context) {
	r.consumer.Start(ctx)
	r.producer.Start(ctx)

	go func() {
		<-ctx.Done()
		r.Close()
	}()
}

func (r *retranslator) Close() {
	r.consumer.Close()
	close(r.events)
	r.producer.Close()
	r.workerPool.StopWait()
}
