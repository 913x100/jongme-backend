package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	UserID      string             `bson:"user_id" json:"user_id"`
	AccessToken string             `bson:"access_token" json:"access_token"`
	Name        string             `bson:"name" json:"name"`
	CreatedOn   time.Time          `bson:"created_on" json:"created_on"`
}

func (u *User) New() *User {
	return &User{
		ID:          primitive.NewObjectID(),
		UserID:      u.UserID,
		AccessToken: u.AccessToken,
		Name:        u.Name,
		CreatedOn:   time.Now(),
	}
}
