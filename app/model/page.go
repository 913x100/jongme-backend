package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// `bson:"_id" json:"id"`
type Page struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	PageID       string             `bson:"page_id" json:"page_id"`
	CreatedOn    time.Time          `bson:"created_on" json:"created_on"`
	OpeningHours []OpeningHours     `bson:"opening_hours" json:"opening_hours"`
}

type OpeningHours struct {
	DayOfWeek int64 `bson:"day" json:"day"`
	Open      int64 `bson:"open" json:"open"`
	Close     int64 `bson:"close" json:"close"`
}

func (p *Page) New() *Page {
	return &Page{
		ID:           primitive.NewObjectID(),
		PageID:       p.PageID,
		CreatedOn:    time.Now(),
		OpeningHours: p.OpeningHours,
	}
}
