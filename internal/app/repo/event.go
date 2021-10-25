package repo

import (
	"github.com/ozonmp/est-water-api/internal/model"
)

type EventRepo interface {
	Lock(n uint64) ([]model.WaterEvent, error)
	Unlock(eventIDs []uint64) error

	Add(event []model.WaterEvent) error
	Remove(eventIDs []uint64) error
}
