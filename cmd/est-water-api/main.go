package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ozonmp/est-water-api/internal/app/retranslator"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)

	cfg := retranslator.Config{
		ChannelSize:   512,

		ConsumerCount: 2,
		ConsumeSize:   10,
		ConsumeTimeout: time.Millisecond,

		ProducerCount: 1,
		WorkerCount:   2,
		WorkerBatchSize: 5,
		WorkerBatchTimeout: time.Millisecond*100,
	}

	transponder := retranslator.NewRetranslator(cfg)
	transponder.Start(ctx)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	<-sigs

	cancel()

	transponder.Close()
}