package migrations

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Migration001_AddInstrumentIndex(db *mongo.Database) error {
	instruments := db.Collection("instruments")

	_, err := instruments.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    map[string]interface{}{"name": 1},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		return err
	}
	log.Println("✅ Migration001_AddInstrumentIndex applied")
	return nil
}
