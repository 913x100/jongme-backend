package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	PageID      string             `bson:"page_id" json:"page_id"`
	UserID      string             `bson:"user_id" json:"user_id"`
	ServiceID   string             `bson:"service_id" json:"service_id"`
	Status      int                `bson:"status" json:"status"` // 0 = waiting 1 = confirmed
	BookingDate string             `bson:"booking_date" json:"booking_date"`
	BookingTime string             `bson:"booking_time" json:"booking_time"`
	Phone       string             `bson:"phone" json:"phone"`
	CreatedOn   time.Time          `bson:"created_on" json:"created_on"`
}

type BookingUser struct {
	Name   string `bson:"name" json:"name"`
	UserID string `bson:"user_id" json:"user_id"`
}

func (b *Booking) New() *Booking {
	return &Booking{
		ID:          primitive.NewObjectID(),
		PageID:      b.PageID,
		UserID:      b.UserID,
		ServiceID:   b.ServiceID,
		BookingDate: b.BookingDate,
		BookingTime: b.BookingTime,
		Status:      b.Status,
		// Time:      b.Time,
		CreatedOn: time.Now(),
	}
}

type AggregateBooking struct {
	ID          string   `bson:"_id" json:"_id"`
	Count       int      `bson:"count" json:"count"`
	PageID      string   `bson:"page_id" json:"page_id"`
	ServiceID   string   `bson:"service_id" json:"service_id"`
	Status      int      `bson:"status" json:"status"`
	BookingDate string   `bson:"booking_date" json:"booking_date"`
	Users       []string `bson:"users" json:"users"`
}
