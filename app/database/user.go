package database

import (
	"context"
	"encoding/json"
	"fmt"
	"jongme/app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var collection string = "users"

func (m *Mongo) CreateUser(user *model.User) error {
	// opts := options
	// 				.FindOneAndReplace()
	// 				.SetUpsert(true)
	// 				.SetSort(bson.D{{"user_id", 1}})
	upsert := true
	opts := &options.FindOneAndReplaceOptions{
		Upsert: &upsert,
		Sort:   bson.D{{"user_id", 1}},
	}
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

func (m *Mongo) GetUser(query interface{}) (*model.User, error) {
	// TODO make code more clean

	// if reflect.TypeOf(t).Kind() != reflect.Struct {
	// 	return nil, errors.New("input type is not a struct.")
	// }

	// query, err := json.Marshal(t)
	// if err != nil {
	// 	// TODO add error
	// 	fmt.Println("error")
	// 	return nil, err
	// }

	var filter bson.M
	var user *model.User

	err := json.Unmarshal(query.([]byte), &filter)

	err = m.DB.Collection("users").FindOne(
		context.Background(),
		bson.M(filter),
	).Decode(&user)

	if err != nil {
		fmt.Println("Errorrrr!")
		return nil, err
	}

	return user, nil

}
