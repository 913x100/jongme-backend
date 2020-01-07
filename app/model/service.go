package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	ID                primitive.ObjectID `bson:"_id" json:"_id"`
	PageID            string             `bson:"page_id" json:"page_id"`
	CreatedOn         time.Time          `bson:"created_on" json:"created_on"`
	Name              string             `bson:"name" json:"name"`
	Type              string             `bson:"type" json:"type"`
	Quantity          int64              `bson:"quantity" json:"quantity"`
	MinimumTimeLength int64              `bson:"minimum_time_length" json:"minimum_time_length"`
	StartTime         int64              `bson:"start_time" json:"start_time"`
	EndTime           int64              `bson:"end_time" json:"end_time"`
}

func (s *Service) New() *Service {
	return &Service{
		ID:                primitive.NewObjectID(),
		PageID:            s.PageID,
		CreatedOn:         time.Now(),
		Name:              s.Name,
		Type:              s.Type,
		Quantity:          s.Quantity,
		MinimumTimeLength: s.MinimumTimeLength,
		StartTime:         s.StartTime,
		EndTime:           s.EndTime,
	}
}
