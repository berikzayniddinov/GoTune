package migrations

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Migration001_AddCartUserIndex(db *mongo.Database) error {
	carts := db.Collection("carts")

	_, err := carts.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    map[string]interface{}{"userId": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}
	log.Println("✅ Migration001_AddCartUserIndex applied")
	return nil
}
