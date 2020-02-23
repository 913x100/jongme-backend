package database

import (
	"context"
	"jongme/app/model"

	"go.mongodb.org/mongo-driver/bson"
)

func (m *Mongo) CreatePage(page *model.Page) (*model.Page, error) {
	_, err := m.DB.Collection("pages").InsertOne(context.Background(), page)

	return page, err

}

func (m *Mongo) GetPages() ([]*model.Page, error) {
	pages := []*model.Page{}

	cursor, err := m.DB.Collection("pages").
		Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		page := &model.Page{}
		if err := cursor.Decode(page); err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}

	return pages, nil
}

func (m *Mongo) GetPageByID(id string) (*model.Page, error) {
	var page *model.Page
	filter := bson.D{{"page_id", id}}

	err := m.DB.Collection("pages").FindOne(context.Background(), filter).Decode(&page)

	return page, err
}

func (m *Mongo) UpdatePage(id string, page *model.Page) (*model.Page, error) {
	doc, err := toDoc(page)
	//check error

	filter := bson.D{{"page_id", page.PageID}}
	update := bson.M{
		"$set": doc,
	}

	_, err = m.DB.Collection("pages").UpdateOne(
		context.Background(),
		filter,
		update,
	)

	return page, err
}
