package api

import (
	"encoding/json"
	"jongme/app/errs"
	"jongme/app/model"
	"strconv"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingDatabase interface {
	CreateBooking(booking *model.Booking) (*model.Booking, error)
	GetBookings() ([]*model.Booking, error)
	GetBookingsAccordingFilter(query []bson.M) ([]*model.Booking, error)
	GetAggregateBookings(query []bson.M) ([]*model.AggregateBooking, error)
	UpdateBooking(booking *model.Booking) (*model.Booking, error)
	DeleteBookingByID(id primitive.ObjectID) error
}

type BookingAPI struct {
	DB BookingDatabase
}

func (b *BookingAPI) CreateBooking(ctx *fasthttp.RequestCtx) error {
	ctx.SetContentType("application/json; charset=utf-8")
	if !ctx.IsPost() {
		return errs.NewHTTPError(nil, 405, "Method not allowed.")
	}
	input := model.Booking{}

	// fmt.Println("book")

	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		return errs.NewHTTPError(err, 400, "Bad request : invalid JSON.")
	}
	// fmt.Println(booking)
	// if err := s.Validate.Struct(service); err != nil {
	// 	return errs.NewHTTPError(err, 400, "Bad request : validation failed.")
	// }

	// service := model.Service{
	// 	PageID: input.PageID,
	// }
	booking := input.New()
	_, err := b.DB.CreateBooking(booking)
	if err != nil {
		return errs.NewHTTPError(err, 500, "Internal server error.")
	}
	json.NewEncoder(ctx).Encode(booking)
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

	id := string(ctx.FormValue("id"))
	pageID := string(ctx.FormValue("page_id"))
	serviceID := string(ctx.FormValue("service_id"))
	userID := string(ctx.FormValue("user_id"))
	name := string(ctx.FormValue("name"))
	status := string(ctx.FormValue("status"))
	year := string(ctx.FormValue("year"))
	month := string(ctx.FormValue("month"))
	day := string(ctx.FormValue("day"))
	time := string(ctx.FormValue("time"))

	filter := []bson.M{}

	if id != "" {
		s, _ := primitive.ObjectIDFromHex(id)
		filter = append(filter, bson.M{"_id": bson.M{"$eq": s}})
	}
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
	if year != "" {
		filter = append(filter, bson.M{"year": bson.M{"$eq": year}})
	}
	if month != "" {
		filter = append(filter, bson.M{"month": bson.M{"$eq": month}})
	}
	if day != "" {
		filter = append(filter, bson.M{"day": bson.M{"$eq": day}})
	}
	if time != "" {
		filter = append(filter, bson.M{"time": bson.M{"$eq": time}})
	}
	if status != "" {
		s, _ := strconv.Atoi(status)
		filter = append(filter, bson.M{"status": bson.M{"$eq": s}})
	}

	bookings, err := b.DB.GetBookingsAccordingFilter(filter)
	if err != nil {
		return errs.NewHTTPError(err, 500, "Internal server error.")
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	json.NewEncoder(ctx).Encode(bookings)
	return nil
}

func (b *BookingAPI) UpdateBookingByID(ctx *fasthttp.RequestCtx) error {
	if !ctx.IsPut() {
		return errs.NewHTTPError(nil, 405, "Method not allowed.")
	}
	// fmt.Println("update")
	input := model.Booking{}

	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		return errs.NewHTTPError(err, 400, "Bad request : invalid JSON.")
	}

	// if err := s.Validate.Struct(input); err != nil {
	// 	return errs.NewHTTPError(err, 400, "Bad request : validation failed.")
	// }
	_, err := withID(ctx, "id")
	if err != nil {
		return errs.NewHTTPError(err, 400, "Bad request: 'invalid id.")
	}

	_, err = b.DB.UpdateBooking(&input)

	if err != nil {
		return errs.NewHTTPError(err, 404, "service down not exists.")
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	return nil
}

func (b *BookingAPI) DeleteBookingByID(ctx *fasthttp.RequestCtx) error {
	if !ctx.IsDelete() {
		return errs.NewHTTPError(nil, 405, "Method not allowed.")
	}

	id, err := withID(ctx, "id")
	if err != nil {
		return errs.NewHTTPError(err, 400, "Bad request: 'invalid id.")
	}

	if err := b.DB.DeleteBookingByID(id); err != nil {
		return errs.NewHTTPError(err, 500, "Internal server error.")
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	return nil
}
