package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gotune/events"
	"gotune/users/internal/config"
	"gotune/users/internal/repository"
	"gotune/users/internal/service"
	"gotune/users/proto"
	"log"
	"net"
)

const (
	mongoURI = "mongodb://localhost:27017"
	dbName   = "gotune_users"
)

func main() {
	mongoClient := config.ConnectMon(mongoURI)
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Ошибка отключения MongoDB: %v", err)
		}
	}()
	db := mongoClient.Database(dbName)
	userRepo := repository.NewUserRepository(db)

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // или поменяй порт, если другой
	})
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Printf("Ошибка закрытия Redis: %v", err)
		}
	}()

	eventPublisher := events.NewEventPublish("amqp://guest:guest@localhost:5672/")

	userService := service.NewUserService(userRepo, eventPublisher, rdb)

	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, userService)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Не удалось слушать порт: %v", err)
	}

	log.Println("Users Service запущен на порту 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка запуска gRPC сервера: %v", err)
	}
}
