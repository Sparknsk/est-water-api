package sender

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/ozonmp/est-water-api/internal/logger"
	"github.com/ozonmp/est-water-api/internal/model"
)

//go:generate mockgen -destination=../../mocks/sender_mock.go -package=mocks github.com/ozonmp/est-water-api/internal/app/sender EventSender
type EventSender interface {
	Send(ctx context.Context, event *model.WaterEvent) error
}

type eventSender struct {
	producer sarama.SyncProducer
	topicName string
}

func NewEventSender(brokersAddr []string, topicName string) (EventSender, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	syncProducer, err := sarama.NewSyncProducer(brokersAddr, config)
	if err != nil {
		return nil, errors.Wrap(err, "sarama.NewSyncProducer() failed")
	}

	return &eventSender{
		producer: syncProducer,
		topicName: topicName,
	}, nil
}

func (es *eventSender) Send(ctx context.Context, event *model.WaterEvent) error {
	logger.DebugKV(ctx, fmt.Sprintf("Send event to broker: %v", event))

	eventPb := event.ModelWaterEventToProtobufWaterEvent()
	message, err := proto.Marshal(eventPb)
	if err != nil {
		return errors.Wrap(err, "proto.Marshal() failed")
	}

	msg := &sarama.ProducerMessage{
		Topic: es.topicName,
		Partition: -1,
		Value: sarama.ByteEncoder(message),
	}

	if _, _, err = es.producer.SendMessage(msg); err != nil {
		return errors.Wrap(err, "kp.producer.SendMessage() failed")
	}

	return nil
}