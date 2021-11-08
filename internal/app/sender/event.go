package sender

import (
	"github.com/ozonmp/est-water-api/internal/model"
)

//go:generate mockgen -destination=../../mocks/sender_mock.go -package=mocks github.com/ozonmp/est-water-api/internal/app/sender EventSender
type EventSender interface {
	Send(event *model.WaterEvent) error
}
