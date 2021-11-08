package model

import "fmt"

const (
	Created EventType = iota
	Updated
	Removed

	Deferred EventStatus = iota
	Processed
)

type EventType uint8

type EventStatus uint8

type WaterEvent struct {
	ID     uint64
	Type   EventType
	Status EventStatus
	Entity *Water
}

type Water struct {
	Id uint64
	Name string
	Model string
	Manufacturer string
	Material string
	Speed uint32
}

func NewWater(id uint64, name string, model string, manufacturer string, material string, speed uint32) *Water {
	return &Water{
		id,
		name,
		model,
		manufacturer,
		material,
		speed,
	}
}

func (a Water) String() string {
	return fmt.Sprintf("id=%d, name=%s, model=%s, manufacturer=%s, material=%s, speed=%d", a.Id, a.Name, a.Model, a.Manufacturer, a.Material, a.Speed)
}