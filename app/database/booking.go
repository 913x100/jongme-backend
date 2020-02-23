package database

import (
	"context"
	"fmt"
	"jongme/app/model"

	"go.mongodb.org/mongo-driver/bson"
)

func (m *Mongo) CreateBooking(booking *model.Booking) error {
	_, err := m.DB.Collection("bookings").InsertOne(context.Background(), booking)
	return err
}

func (m *Mongo) GetBookings() ([]*model.Booking, error) {
	bookings := []*model.Booking{}

	cursor, err := m.DB.Collection("bookings").
		Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		booking := &model.Booking{}
		if err := cursor.Decode(booking); err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (m *Mongo) GetBookingsAccordingFilter(query []bson.M) ([]*model.Booking, error) {
	fmt.Println("Get Book")

	filter := bson.M{}
	if query != nil {
		if len(query) > 0 {
			filter = bson.M{"$and": query}
		}
	}
	// fmt.Println(filter)
	data, err := m.DB.Collection("bookings").Find(
		context.Background(),
		filter,
		nil)
	if err != nil {
		return nil, err
	}
	defer data.Close(context.Background())

	var result []*model.Booking
	for data.Next(context.Background()) {
		l := &model.Booking{}
		err = data.Decode(&l)
		if err != nil {
			return nil, err
		}
		result = append(result, l)
	}
	return result, nil
}

func (m *Mongo) GetAggregateBookings(query []bson.M) ([]*model.AggregateBooking, error) {

	fmt.Println(query)
	data, err := m.DB.Collection("bookings").Aggregate(context.Background(), query)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer data.Close(context.Background())

	var result []*model.AggregateBooking
	for data.Next(context.Background()) {
		fmt.Println(data)
		l := &model.AggregateBooking{}
		err = data.Decode(&l)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		result = append(result, l)
	}
	return result, nil
}
