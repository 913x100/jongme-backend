package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Page struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	PageID      string             `bson:"page_id" json:"page_id"`
	Name        string             `bson:"name" json:"name"`
	AccessToken string             `bson:"access_token" json:"access_token"`
	PageHours   []PageHours        `bson:"page_hours" json:"page_hours"`
	CreatedOn   time.Time          `bson:"created_on" json:"created_on"`
}

type PageHours struct {
	DayOfWeek  int64 `bson:"day" json:"day"`
	Open       int64 `bson:"open" json:"open"`
	Close      int64 `bson:"close" json:"close"`
	IsBreak    bool  `bson:"is_break" json:"is_break"`
	BreakStart int64 `bson:"break_start" json:"break_start"`
	BreakEnd   int64 `bson:"break_end" json:"break_end"`
}

type PageBooking struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	PageID        string             `bson:"page_id" json:"page_id"`
	BookingDays   int                `bson:"booking_days" json:"booking_days"`
	IsAutoConfirm bool               `bson:"is_auto_confirm" json:"is_auto_confirm"`
	IsAskingPhone bool               `bson:"is_asking_phone" json:"is_asking_phone"`
	IsAskingName  bool               `bson:"is_asking_name" json:"is_asking_name"`
}

func (p *Page) New() *Page {
	return &Page{
		ID:          primitive.NewObjectID(),
		PageID:      p.PageID,
		Name:        p.Name,
		AccessToken: p.AccessToken,
		PageHours:   p.PageHours,
		CreatedOn:   time.Now(),
	}
}

func (p *PageBooking) New() *PageBooking {
	return &PageBooking{
		ID:            primitive.NewObjectID(),
		PageID:        p.PageID,
		BookingDays:   p.BookingDays,
		IsAutoConfirm: p.IsAutoConfirm,
		IsAskingPhone: p.IsAskingPhone,
		IsAskingName:  p.IsAskingName,
	}
}
