package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Page struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	PageID      string             `bson:"page_id" json:"page_id,omitempty"`
	Name        string             `bson:"name" json:"name,omitempty"`
	AccessToken string             `bson:"access_token" json:"access_token,omitempty"`
	StartTime   string             `bson:"start_time" json:"start_time"`
	EndTime     string             `bson:"end_time" json:"end_time"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	IsBreak     bool               `bson:"is_break" json:"is_break"`
	BreakStart  string             `bson:"break_start" json:"break_start"`
	BreakEnd    string             `bson:"break_end" json:"break_end"`
	Sun         bool               `bson:"sun" json:"sun"`
	Mon         bool               `bson:"mon" json:"mon"`
	Tue         bool               `bson:"tue" json:"tue"`
	Wed         bool               `bson:"wed" json:"wed"`
	Thu         bool               `bson:"thu" json:"thu"`
	Fri         bool               `bson:"fri" json:"fri"`
	Sat         bool               `bson:"sat" json:"sat"`
	UpdatedOn   time.Time          `bson:"updated_on" json:"updated_on,omitempty"`
	CreatedOn   time.Time          `bson:"created_on" json:"created_on,omitempty"`
}

type DayOfWeek struct {
}

type PageHours struct {
	DayOfWeek  int8  `bson:"day" json:"day"`
	Open       int64 `bson:"open" json:"open"`
	Close      int64 `bson:"close" json:"close"`
	IsBreak    bool  `bson:"is_break" json:"is_break"`
	BreakStart int64 `bson:"break_start" json:"break_start"`
	BreakEnd   int64 `bson:"break_end" json:"break_end"`
}

type UpdatePageToken struct {
	PageID      string `bson:"page_id" json:"page_id"`
	AccessToken string `bson:"access_token" json:"access_token"`
	Name        string `bson:"name" json:"name"`
	// UpdatedOn   time.Time `bson:"updated_on" json:"updated_on,omitempty"`
	// CreatedOn   time.Time `bson:"created_on" json:"created_on,omitempty"`
}

type UpdatePage struct {
	PageID     string `bson:"page_id" json:"page_id,omitempty"`
	StartTime  string `bson:"start_time" json:"start_time"`
	EndTime    string `bson:"end_time" json:"end_time"`
	IsActive   bool   `bson:"is_active" json:"is_active"`
	IsBreak    bool   `bson:"is_break" json:"is_break"`
	Sun        bool   `bson:"sun" json:"sun"`
	Mon        bool   `bson:"mon" json:"mon"`
	Tue        bool   `bson:"tue" json:"tue"`
	Wed        bool   `bson:"wed" json:"wed"`
	Thu        bool   `bson:"thu" json:"thu"`
	Fri        bool   `bson:"fri" json:"fri"`
	Sat        bool   `bson:"sat" json:"sat"`
	BreakStart string `bson:"break_start" json:"break_start"`
	BreakEnd   string `bson:"break_end" json:"break_end"`
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
		// PageHours:   p.PageHours,
		StartTime:  p.StartTime,
		EndTime:    p.EndTime,
		IsActive:   p.IsActive,
		IsBreak:    p.IsBreak,
		BreakStart: p.BreakStart,
		BreakEnd:   p.BreakEnd,
		Sun:        p.Sun,
		Mon:        p.Mon,
		Tue:        p.Tue,
		Wed:        p.Wed,
		Thu:        p.Thu,
		Fri:        p.Fri,
		Sat:        p.Sat,
		UpdatedOn:  p.UpdatedOn,
		CreatedOn:  time.Now(),
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
