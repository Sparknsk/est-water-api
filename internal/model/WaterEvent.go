package model

import (
	"database/sql/driver"
	"errors"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/ozonmp/est-water-api/pkg/est-water-api"
)

const (
	Created EventType = iota
	Removed
	UpdatedName
	UpdatedModel
	UpdatedManufacturer
	UpdatedMaterial
	UpdatedSpeed
)

const (
	Locked EventStatus = iota
	Unlocked
)

type EventType uint8

type EventStatus uint8

type WaterEvent struct {
	ID uint64 `db:"id"`
	WaterId uint64 `db:"water_id"`
	Type EventType `db:"type"`
	Status EventStatus `db:"status"`
	Entity *Water `db:"payload"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func (we *WaterEvent) ModelWaterEventToProtobufWaterEvent() *pb.WaterEvent {
	return &pb.WaterEvent{
		Id: we.ID,
		WaterId: we.WaterId,
		Type: pb.WaterEvent_Type(we.Type),
		Status: pb.WaterEvent_Status(we.Status),
		Entity: we.Entity.ModelWaterToProtobufWater(),
		CreatedAt: timestamppb.New(*we.CreatedAt),
		UpdatedAt: timestamppb.New(*we.UpdatedAt),
	}
}

func (et EventType) Value() (driver.Value, error) {
	var eventType string
	switch et {
	case Created:
		eventType = "created"
	case Removed:
		eventType = "removed"
	case UpdatedName:
		eventType = "updated_name"
	case UpdatedModel:
		eventType = "updated_model"
	case UpdatedMaterial:
		eventType = "updated_material"
	case UpdatedManufacturer:
		eventType = "updated_manufacturer"
	case UpdatedSpeed:
		eventType = "updated_speed"
	default:
		return nil, errors.New("undefined event type")
	}
	return eventType, nil
}

func (et *EventType) Scan(src interface{}) (err error) {
	if src == nil {
		return nil
	}
	var eventType EventType
	switch src {
	case "created":
		eventType = Created
	case "removed":
		eventType = Removed
	case "updated_name":
		eventType = UpdatedName
	case "updated_model":
		eventType = UpdatedModel
	case "updated_material":
		eventType = UpdatedMaterial
	case "updated_manufacturer":
		eventType = UpdatedManufacturer
	case "updated_speed":
		eventType = UpdatedSpeed
	default:
		return errors.New("undefined event type")
	}

	if err != nil {
		return err
	}

	*et = eventType
	return nil
}

func (es EventStatus) Value() (driver.Value, error) {
	var eventStatus string
	switch es {
	case Unlocked:
		eventStatus = "unlock"
	case Locked:
		eventStatus = "lock"
	default:
		return nil, errors.New("undefined event status")
	}
	return eventStatus, nil
}

func (es *EventStatus) Scan(src interface{}) (err error) {
	if src == nil {
		return nil
	}
	var eventStatus EventStatus
	switch src {
	case "unlock":
		eventStatus = Unlocked
	case "lock":
		eventStatus = Locked
	default:
		return errors.New("undefined event status")
	}

	if err != nil {
		return err
	}

	*es = eventStatus
	return nil
}