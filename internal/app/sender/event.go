package sender

import (
	"fmt"
	"github.com/ozonmp/est-water-api/internal/model"
	"time"
)

//go:generate mockgen -destination=../../mocks/sender_mock.go -package=mocks github.com/ozonmp/est-water-api/internal/app/sender EventSender
type EventSender interface {
	Send(event *model.WaterEvent) error
}

type eventSender struct {

}

func NewEventSender() EventSender {
	return &eventSender{}
}

func (* eventSender) Send(event *model.WaterEvent) error {
	time.Sleep(time.Millisecond*100)
	fmt.Printf("SEND EVENT: %v\n", event)
	return nil
}
