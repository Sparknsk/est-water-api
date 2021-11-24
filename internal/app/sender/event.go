package sender

import (
	"context"
	"fmt"
	"time"

	"github.com/ozonmp/est-water-api/internal/logger"
	"github.com/ozonmp/est-water-api/internal/model"
)

//go:generate mockgen -destination=../../mocks/sender_mock.go -package=mocks github.com/ozonmp/est-water-api/internal/app/sender EventSender
type EventSender interface {
	Send(ctx context.Context, event *model.WaterEvent) error
}

type eventSender struct {

}

func NewEventSender() EventSender {
	return &eventSender{}
}

func (* eventSender) Send(ctx context.Context, event *model.WaterEvent) error {
	time.Sleep(time.Millisecond*100)
	logger.DebugKV(ctx, fmt.Sprintf("Send event: %v", event))
	return nil
}