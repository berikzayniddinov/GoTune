// users/migrations/migration_001_add_user_index.go
package migrations

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Migration001_AddUserIndex(db *mongo.Database) error {
	users := db.Collection("users")

	_, err := users.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    map[string]interface{}{"email": 1},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		return err
	}
	log.Println("✅ Migration001_AddUserIndex applied")
	return nil
}
