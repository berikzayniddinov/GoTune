package migrations

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Migration001_AddOrderUserIndex(db *mongo.Database) error {
	orders := db.Collection("orders")

	_, err := orders.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    map[string]interface{}{"userId": 1},
		Options: options.Index().SetUnique(false), // для заказов уникальность необязательна
	})
	if err != nil {
		return err
	}
	log.Println("✅ Migration001_AddOrderUserIndex applied")
	return nil
}
