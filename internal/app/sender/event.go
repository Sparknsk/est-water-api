package sender

import (
	"github.com/ozonmp/est-water-api/internal/model"
)

type EventSender interface {
	Send(event *model.WaterEvent) error
}
