package model

import (
	"database/sql/driver"
	"errors"
	"time"
)

const (
	Created EventType = iota
	Updated
	Removed

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

func (et EventType) Value() (driver.Value, error) {
	var eventType string
	switch et {
	case Created:
		eventType = "created"
	case Updated:
		eventType = "updated"
	case Removed:
		eventType = "removed"
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
	case "updated":
		eventType = Updated
	case "removed":
		eventType = Removed
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