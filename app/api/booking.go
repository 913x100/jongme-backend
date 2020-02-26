package api

import (
	"encoding/json"
	"fmt"
	"jongme/app/errs"
	"jongme/app/model"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingDatabase interface {
	CreateBooking(booking *model.Booking) error
	GetBookings() ([]*model.Booking, error)
	GetBookingsAccordingFilter(query []bson.M) ([]*model.Booking, error)
	GetAggregateBookings(query []bson.M) ([]*model.AggregateBooking, error)
}

type BookingAPI struct {
	DB BookingDatabase
}

func (b *BookingAPI) CreateBooking(ctx *fasthttp.RequestCtx) error {
	ctx.SetContentType("application/json; charset=utf-8")
	if !ctx.IsPost() {
		return errs.NewHTTPError(nil, 405, "Method not allowed.")
	}
	booking := model.Booking{}

	fmt.Println("book")

	if err := json.Unmarshal(ctx.PostBody(), &booking); err != nil {
		return errs.NewHTTPError(err, 400, "Bad request : invalid JSON.")
	}
	fmt.Println(booking)
	// if err := s.Validate.Struct(service); err != nil {
	// 	return errs.NewHTTPError(err, 400, "Bad request : validation failed.")
	// }

	// service := model.Service{
	// 	PageID: input.PageID,
	// }

	err := b.DB.CreateBooking(booking.New())
	fmt.Println(err)
	if err != nil {
		return errs.NewHTTPError(err, 500, "Internal server error.")
	}
	ctx.SetStatusCode(fasthttp.StatusCreated)
	return nil
}

func (b *BookingAPI) GetBookingByService(ctx *fasthttp.RequestCtx) error {
	ctx.SetContentType("application/json; charset=utf-8")
	if !ctx.IsGet() {
		return errs.NewHTTPError(nil, 405, "Method not allowed.")
	}

	serviceID := string(ctx.FormValue("service_id"))
	year := string(ctx.FormValue("year"))
	month := string(ctx.FormValue("month"))
	day := string(ctx.FormValue("day"))

	query := []bson.M{
		{"$match": bson.D{
			{"service_id", serviceID},
			{"year", year},
			{"month", month},
			{"day", day},
		}},
		{"$group": bson.M{"_id": "$time",
			"count":      bson.M{"$sum": 1},
			"page_id":    bson.M{"$first": "$page_id"},
			"service_id": bson.M{"$first": "$service_id"},
			"year":       bson.M{"$first": "$year"},
			"month":      bson.M{"$first": "$month"},
			"day":        bson.M{"$first": "$day"},
			"status":     bson.M{"$first": "$status"},
			"users":      bson.M{"$addToSet": "$user_id"},
		}},
	}

	bookings, err := b.DB.GetAggregateBookings(query)
	if err != nil {
		return err
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	json.NewEncoder(ctx).Encode(bookings)

	return nil
}

func (b *BookingAPI) GetBookings(ctx *fasthttp.RequestCtx) error {
	ctx.SetContentType("application/json;charset=utf-8")

	if !ctx.IsGet() {
		return errs.NewHTTPError(nil, 405, "Method not allowed.")
	}

	bookings, err := b.DB.GetBookings()

	if err != nil {
		return errs.NewHTTPError(err, 500, "Internal server error.")
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	json.NewEncoder(ctx).Encode(bookings)
	return nil
}

func (b *BookingAPI) GetBookingByFilter(ctx *fasthttp.RequestCtx) error {
	ctx.SetContentType("application/json;charset=utf-8")

	pageID := string(ctx.FormValue("page_id"))
	serviceID := string(ctx.FormValue("service_id"))
	userID := string(ctx.FormValue("user_id"))
	name := string(ctx.FormValue("name"))
	bookingTime := string(ctx.FormValue("booking_time"))

	filter := []bson.M{}

	if pageID != "" {
		filter = append(filter, bson.M{"page_id": bson.M{"$eq": pageID}})
	}
	if serviceID != "" {
		filter = append(filter, bson.M{"service_id": bson.M{"$eq": serviceID}})
	}
	if userID != "" {
		filter = append(filter, bson.M{"user_id": bson.M{"$eq": userID}})
	}
	if name != "" {
		filter = append(filter, bson.M{"name": bson.M{"$eq": name}})
	}
	if bookingTime != "" {
		filter = append(filter, bson.M{"booking_time": bson.M{"$eq": bookingTime}})
	}

	bookings, err := b.DB.GetBookingsAccordingFilter(filter)
	if err != nil {
		return errs.NewHTTPError(err, 500, "Internal server error.")
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	json.NewEncoder(ctx).Encode(bookings)
	return nil
}
