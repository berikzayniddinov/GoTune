package main

import (
	"context"
	"google.golang.org/grpc/reflection"
	"log"
	"net"

	"gotune/users/internal/config"
	"gotune/users/internal/repository"
	"gotune/users/internal/service"
	"gotune/users/proto"

	"google.golang.org/grpc"
)

const (
	mongoURI = "mongodb://localhost:27017"
	dbName   = "gotune_users"
)

func main() {
	mongoClient := config.ConnectMongo(mongoURI)
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Ошибка отключения MongoDB: %v", err)
		}
	}()

	db := mongoClient.Database(dbName)
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

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
