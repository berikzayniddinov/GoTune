package main

import (
	"context"
	"gotune/events"
	"log"
	"net"

	"gotune/order/internal/config"
	"gotune/order/internal/repository"
	"gotune/order/internal/service"
	"gotune/order/proto"
	usersproto "gotune/users/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	mongoURI           = "mongodb://localhost:27017"
	dbName             = "gotune_order"
	userServiceAddress = "localhost:50051"
)

func main() {
	mongoClient := config.ConnectMongo(mongoURI)
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Ошибка отключения MongoDB: %v", err)
		}
	}()

	db := mongoClient.Database(dbName)

	orderRepo := repository.NewOrderRepository(db)
	eventPublisher := events.NewEventPublish("amqp://guest:guest@localhost:5672/")
	userConn, err := grpc.Dial(userServiceAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не удалось подключиться к UserService: %v", err)
	}
	defer userConn.Close()

	userClient := usersproto.NewUserServiceClient(userConn)

	orderService := service.NewOrderService(orderRepo, userClient, eventPublisher)

	grpcServer := grpc.NewServer()

	proto.RegisterOrderServiceServer(grpcServer, orderService)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("Не удалось слушать порт: %v", err)
	}

	log.Println("Order Service запущен на порту 50054")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка запуска gRPC сервера: %v", err)
	}
}
