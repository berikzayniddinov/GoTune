package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMon(uri string) *mongo.Client {
	clientOpts := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("Ошибка подключения к MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Ошибка при проверке подключения к MongoDB: %v", err)
	}

	log.Println("Успешно подключено к MongoDB")
	return client
}
