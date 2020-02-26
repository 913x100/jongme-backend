package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	ID                primitive.ObjectID `bson:"_id" json:"_id"`
	PageID            string             `bson:"page_id" json:"page_id"`
	IsActive          bool               `bson:"is_active" json:"is_active"`
	Name              string             `bson:"name" json:"name"`
	ImageUrl          string             `bson:"image_url" json:"image_url"`
	UnitType          string             `bson:"unit_type" json:"unit_type"`
	UnitQuantity      int                `bson:"unit_quantity" json:"unit_quantity"`
	MinimumTimeLength int                `bson:"minimum_time_length" json:"minimum_time_length"`
	IsTimeAdjust      bool               `bson:"is_time_adjust" json:"is_time_adjust"`
	StartTime         string             `bson:"start_time" json:"start_time"`
	EndTime           string             `bson:"end_time" json:"end_time"`
	CreatedOn         time.Time          `bson:"created_on" json:"created_on"`
}

func (s *Service) New() *Service {
	return &Service{
		ID:                primitive.NewObjectID(),
		PageID:            s.PageID,
		IsActive:          s.IsActive,
		Name:              s.Name,
		ImageUrl:          s.ImageUrl,
		UnitType:          s.UnitType,
		UnitQuantity:      s.UnitQuantity,
		MinimumTimeLength: s.MinimumTimeLength,
		IsTimeAdjust:      s.IsTimeAdjust,
		StartTime:         s.StartTime,
		EndTime:           s.EndTime,
		CreatedOn:         time.Now(),
	}
}

type AggregateService struct {
	Qty      int64      `json:"qty"`
	Services []*Service `json:"services"`
}
