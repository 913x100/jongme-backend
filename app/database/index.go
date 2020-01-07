package database

import (
	"context"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// CreateUniqueIndex create UniqueIndex
func (m *Mongo) CreateUniqueIndex(collection string, keys ...string) {
	// Composite index
	keysDoc := bsonx.Doc{}
	for _, key := range keys {
		if strings.HasPrefix(key, "-") {
			keysDoc = keysDoc.Append(strings.TrimLeft(key, "-"), bsonx.Int32(-1))
		} else {
			keysDoc = keysDoc.Append(key, bsonx.Int32(1))
		}
	}
	// Create index
	idxRet, err := m.DB.Collection(collection).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    keysDoc,
			Options: options.Index().SetUnique(true),
		},
		options.CreateIndexes().SetMaxTime(10*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("collection.Indexes().CreateOne:", idxRet)
}
