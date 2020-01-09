package database

import (
	"context"
	"fmt"
	"jongme/app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var collection string = "services"

func (m *Mongo) CreateService(service *model.Service) error {
	_, err := m.DB.Collection("services").InsertOne(context.Background(), service)
	return err
}

func (m *Mongo) GetServices(paging *model.Paging) ([]*model.Service, error) {
	services := []*model.Service{}

	cursor, err := m.DB.Collection("services").
		Find(context.Background(), bson.D{},
			&options.FindOptions{
				Skip:  paging.Skip,
				Sort:  bson.D{bson.E{Key: paging.SortKey, Value: paging.SortVal}},
				Limit: paging.Limit,
			})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		service := &model.Service{}
		if err := cursor.Decode(service); err != nil {
			return nil, err
		}
		services = append(services, service)
	}

	return services, nil
}

func (m *Mongo) UpdateService(service *model.Service) error {
	opts := options.FindOneAndReplace()
	filter := bson.D{{"_id", service.ID}}
	fmt.Println(service)
	result := m.DB.Collection("services").
		FindOneAndReplace(context.Background(),
			filter,
			service,
			opts,
		)

	service = new(model.Service)
	if err := result.Decode(service); err != nil {
		return err
	}
	return nil
}

func (m *Mongo) DeleteServiceByID(id primitive.ObjectID) error {
	_, err := m.DB.Collection("services").DeleteOne(context.Background(), bson.D{{Key: "_id", Value: id}})

	return err
}
