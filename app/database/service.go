package database

import (
	"context"
	"jongme/app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var collection string = "services"

func (m *Mongo) CreateService(service *model.Service) (*model.Service, error) {
	_, err := m.DB.Collection("services").InsertOne(context.Background(), service)
	if err != nil {
		return nil, err
	}
	return service, err
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

func (m *Mongo) GetServiceByID(id primitive.ObjectID) (*model.Service, error) {
	var service *model.Service
	filter := bson.D{{"_id", id}}

	err := m.DB.Collection("services").FindOne(context.Background(), filter).Decode(&service)

	return service, err
}

func (m *Mongo) UpdateService(service *model.Service) (*model.Service, error) {
	doc, err := toDoc(service)
	//check error

	filter := bson.D{{"_id", service.ID}}
	update := bson.M{
		"$set": doc,
	}

	_, err = m.DB.Collection("services").UpdateOne(
		context.Background(),
		filter,
		update,
	)

	return service, err

}

func (m *Mongo) DeleteServiceByID(id primitive.ObjectID) error {
	_, err := m.DB.Collection("services").DeleteOne(context.Background(), bson.D{{Key: "_id", Value: id}})

	return err
}

func (m *Mongo) GetServicesAccordingFilter(query []bson.M) ([]*model.Service, error) {

	filter := bson.M{}
	if query != nil {
		if len(query) > 0 {
			filter = bson.M{"$and": query}
		}
	}
	data, err := m.DB.Collection("services").Find(
		context.Background(),
		filter,
		nil)
	if err != nil {
		return nil, err
	}
	defer data.Close(context.Background())

	var result []*model.Service
	for data.Next(context.Background()) {
		l := &model.Service{}
		err = data.Decode(&l)
		if err != nil {
			return nil, err
		}
		result = append(result, l)
	}
	return result, nil
}
