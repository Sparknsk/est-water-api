package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ozonmp/est-water-api/internal/app/retranslator"
)

func main() {

	sigs := make(chan os.Signal, 1)

	cfg := retranslator.Config{
		ChannelSize:   512,

		ConsumerCount: 2,
		ConsumeSize:   10,
		ConsumeTimeout: 10,

		ProducerCount: 28,
		WorkerCount:   2,
	}

	transponder := retranslator.NewRetranslator(cfg)
	transponder.Start()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}