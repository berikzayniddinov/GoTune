package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gotune/events"
	"log"
	"net"

	"gotune/instruments/internal/config"
	"gotune/instruments/internal/repository"
	"gotune/instruments/internal/service"
	"gotune/instruments/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	mongoURI = "mongodb://localhost:27017"
	dbName   = "gotune_instruments"
)

func main() {
	mongoClient := config.ConnectMongo(mongoURI)
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Ошибка отключения MongoDB: %v", err)
		}
	}()

	db := mongoClient.Database(dbName)
	instrumentRepo := repository.NewInstrumentRepositories(db)
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // или поменяй порт, если другой
	})
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Printf("Ошибка закрытия Redis: %v", err)
		}
	}()
	eventPublisher := events.NewEventPublish("amqp://guest:guest@localhost:5672/")
	instrumentService := service.NewInstrumentService(instrumentRepo, eventPublisher, rdb)

	grpcServer := grpc.NewServer()

	proto.RegisterInstrumentServiceServer(grpcServer, instrumentService)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Не удалось слушать порт: %v", err)
	}

	log.Println("Instrument Service запущен на порту 50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка запуска gRPC сервера: %v", err)
	}
}
