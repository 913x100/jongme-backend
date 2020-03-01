package database

import (
	"context"
	"fmt"
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

func (m *Mongo) UpdatePage(p interface{}) (*model.Page, error) {

	var filter bson.D
	var doc *bson.M

	switch p.(type) {
	case *model.UpdatePage:
		page := p.(*model.UpdatePage)
		filter = bson.D{{"page_id", page.PageID}}
		doc, _ = toDoc(page)
		// fmt.Println(doc)
		break
	case *model.UpdatePageToken:
		page := p.(*model.UpdatePageToken)
		filter = bson.D{{"page_id", page.PageID}}
		doc, _ = toDoc(page)
		break
	}
	update := bson.M{
		"$set": doc,
	}

	_, err := m.DB.Collection("pages").UpdateOne(
		context.Background(),
		filter,
		update,
	)

	fmt.Println(err)

	return nil, err
}
