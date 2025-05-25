package migrations

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Migration struct {
	Name string
	Func func(*mongo.Database) error
}

func RunM(db *mongo.Database) error {
	migrations := []Migration{
		{Name: "Migration001_AddCartUserIndex", Func: Migration001_AddCartUserIndex},
	}

	applied := db.Collection("migrations")

	for _, m := range migrations {
		var result bson.M
		err := applied.FindOne(context.Background(), bson.M{"name": m.Name}).Decode(&result)
		if err == nil {
			log.Printf("✅ Миграция %s уже применена, пропускаем", m.Name)
			continue
		}

		if err := m.Func(db); err != nil {
			return err
		}

		_, err = applied.InsertOne(context.Background(), bson.M{
			"name":      m.Name,
			"appliedAt": time.Now(),
		})
		if err != nil {
			return err
		}
		log.Printf("✅ Миграция %s успешно применена", m.Name)
	}

	return nil
}
