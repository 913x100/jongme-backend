package database

import (
	"context"
	"jongme/app/model"

	"go.mongodb.org/mongo-driver/bson"
)

// var collection string = "users"

func (m *Mongo) CreateUser(user *model.User) error {
	filter := bson.D{{"user_id", user.UserID}}
	var queryUser *model.User
	_ = m.DB.Collection("users").FindOne(context.Background(), filter).Decode(&queryUser)

	if queryUser == nil {
		_, err := m.DB.Collection("users").InsertOne(context.Background(), user)
		return err
	}

	update := bson.M{
		"$set": bson.D{
			{"name", user.Name},
			{"updated_on", user.UpdatedOn},
			{"access_token", user.AccessToken},
		},
	}
	_, err := m.DB.Collection("users").UpdateOne(
		context.Background(),
		filter,
		update,
	)

	return err

}

func (m *Mongo) GetUsers() ([]*model.User, error) {
	users := []*model.User{}

	cursor, err := m.DB.Collection("users").
		Find(context.Background(), bson.D{})
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

func (m *Mongo) GetUserByID(id string) (*model.User, error) {
	var user *model.User
	filter := bson.D{{"user_id", id}}

	err := m.DB.Collection("users").FindOne(context.Background(), filter).Decode(&user)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *Mongo) GetUserByUserID(id string) (*model.User, error) {
	var user *model.User
	filter := bson.D{{"user_id", id}}

	err := m.DB.Collection("users").FindOne(context.Background(), filter).Decode(&user)

	if err != nil {
		return nil, err
	}
	return user, nil
}
