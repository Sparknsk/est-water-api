package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/ozonmp/est-water-api/pkg/est-water-api"
)

type Water struct {
	Id uint64 `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Model string `db:"model" json:"model"`
	Manufacturer string `db:"manufacturer" json:"manufacturer"`
	Material string `db:"material" json:"material"`
	Speed uint32 `db:"speed" json:"speed"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
	DeleteStatus bool `db:"delete_status" json:"delete_status"`
}

func NewWater(id uint64, name string, model string, manufacturer string, material string, speed uint32, createdAt *time.Time, updatedAt *time.Time, deleteStatus bool) *Water {
	return &Water{
		id,
		name,
		model,
		manufacturer,
		material,
		speed,
		createdAt,
		updatedAt,
		deleteStatus,
	}
}

func (w Water) String() string {
	return fmt.Sprintf("id=%d, name=%s, model=%s, manufacturer=%s, material=%s, speed=%d, created_at=%s", w.Id, w.Name, w.Model, w.Manufacturer, w.Material, w.Speed, w.CreatedAt)
}

func (w Water) Value() (driver.Value, error) {
	if w == (Water{}) {
		return []byte{}, nil
	}
	return json.Marshal(w)
}

func (w *Water) Scan(src interface{}) (err error) {
	if src == nil {
		return nil
	}
	var water Water
	switch src.(type) {
	case string:
		err = json.Unmarshal([]byte(src.(string)), &water)
	case []byte:
		err = json.Unmarshal(src.([]byte), &water)
	default:
		return errors.New("incompatible type")
	}

	if err != nil {
		return err
	}

	*w = water
	return nil
}

func (w *Water) ModelWaterToProtobufWater() *pb.Water {
	return &pb.Water{
		Id: w.Id,
		Name: w.Name,
		Model: w.Model,
		Manufacturer: w.Manufacturer,
		Material: w.Material,
		Speed: w.Speed,
		CreatedAt: timestamppb.New(*w.CreatedAt),
	}
}