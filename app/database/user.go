package database

import (
	"context"
	"jongme/app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var collection string = "users"

func (m *Mongo) CreateUser(user *model.User) error {
	opts := options.FindOneAndReplace().SetUpsert(true)
	filter := bson.D{{"user_id", user.UserID}}

	result := m.DB.Collection("users").
		FindOneAndReplace(context.Background(),
			filter,
			user,
			opts,
		)
	user = new(model.User)
	if err := result.Decode(user); err != nil {
		return err
	}
	return nil
}

func (m *Mongo) GetUsers(paging *model.Paging) ([]*model.User, error) {
	users := []*model.User{}

	cursor, err := m.DB.Collection("users").
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
		user := &model.User{}
		if err := cursor.Decode(user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
