package main

import (
	"context"
	"log"
	"net"

	"gotune/cart/internal/config"
	"gotune/cart/internal/repository"
	"gotune/cart/internal/service"
	"gotune/cart/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	mongoURI = "mongodb://localhost:27017"
	dbName   = "gotune_cart"
)

func main() {
	mongoClient := config.ConnectMongo(mongoURI)
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Ошибка отключения MongoDB: %v", err)
		}
	}()

	db := mongoClient.Database(dbName)
	cartRepo := repository.NewCartRepositories(db)
	cartService := service.NewCartService(cartRepo)

	grpcServer := grpc.NewServer()

	proto.RegisterCartServiceServer(grpcServer, cartService)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Не удалось слушать порт: %v", err)
	}

	log.Println("Cart Service запущен на порту 50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка запуска gRPC сервера: %v", err)
	}
}
