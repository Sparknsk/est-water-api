package water_service

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/ozonmp/est-water-api/internal/model"
)

func (s *waterService) UpdateWater(
	ctx context.Context,
	waterId uint64,
	waterName string,
	waterModel string,
	waterManufacturer string,
	waterMaterial string,
	waterSpeed uint32,
) (*model.Water, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "waterService.CreateWater()")
	defer span.Finish()
	span.LogKV(
		"event", "service update water",
		"waterId", waterId,
		"waterName", waterName,
		"waterSpeed", waterSpeed,
	)

	water, err := s.waterRepository.Get(ctx, waterId)
	if err != nil {
		return nil, errors.Wrapf(err, "waterRepository.Get() failed with id=%d", waterId)
	}

	if water == nil {
		return nil, WaterNotFound
	}

	ts := time.Now().UTC()
	water.UpdatedAt = &ts

	var waterEvents []model.WaterEvent

	if waterName != water.Name {
		waterEvents = append(
			waterEvents,
			s.updateField(water, "Name", waterName),
		)
	}

	if waterModel != water.Model {
		waterEvents = append(
			waterEvents,
			s.updateField(water, "Model", waterModel),
		)
	}

	if waterManufacturer != water.Manufacturer {
		waterEvents = append(
			waterEvents,
			s.updateField(water, "Manufacturer", waterManufacturer),
		)
	}

	if waterMaterial != water.Material {
		waterEvents = append(
			waterEvents,
			s.updateField(water, "Material", waterMaterial),
		)
	}

	if waterSpeed != water.Speed {
		waterEvents = append(
			waterEvents,
			s.updateField(water, "Speed", waterSpeed),
		)
	}

	if len(waterEvents) > 0 {
		tx, err := s.db.BeginTxx(ctx, nil)
		if err != nil {
			return nil, errors.Wrap(err, "db.BeginTxx() failed")
		}

		if err := s.waterRepository.Update(ctx, water); err != nil {
			return nil, errors.Wrapf(err, "waterRepository.Update() failed with %v", water)
		}

		if err := s.waterEventRepository.Add(ctx, waterEvents); err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, errors.Wrap(err, "tx.Rollback() failed")
			}
			return nil, errors.Wrapf(err,"waterEventRepository.Add() failed with %v", waterEvents)
		}

		if err := tx.Commit(); err != nil {
			return nil, errors.Wrap(err, "tx.Commit() failed")
		}
	}

	return water, nil
}

func (s *waterService) updateField(water *model.Water, field string, value interface{}) model.WaterEvent {
	var eventType model.EventType

	if field == "Name" {
		water.Name = value.(string)
		eventType = model.UpdatedName
	} else if field == "Model" {
		water.Model = value.(string)
		eventType = model.UpdatedModel
	} else if field == "Manufacturer" {
		water.Manufacturer = value.(string)
		eventType = model.UpdatedManufacturer
	} else if field == "Material" {
		water.Material = value.(string)
		eventType = model.UpdatedMaterial
	} else if field == "Speed" {
		water.Speed = value.(uint32)
		eventType = model.UpdatedSpeed
	}

	waterCopy := *water
	return model.WaterEvent{
		WaterId: water.Id,
		Type: eventType,
		Status: model.Unlocked,
		Entity: &waterCopy,
		CreatedAt: waterCopy.UpdatedAt,
	}
}
