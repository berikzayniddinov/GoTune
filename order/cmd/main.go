package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"

	"gotune/events"
	"gotune/order/metrics"
	"gotune/order/migrations"

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
	// 📌 Запуск HTTP-сервера для метрик
	go func() {
		metrics.Register() // Регистрация метрик
		http.Handle("/metrics", promhttp.Handler())
		log.Println("📊 Метрики доступны на :2113/metrics")
		if err := http.ListenAndServe(":2113", nil); err != nil {
			log.Fatalf("Ошибка запуска HTTP-сервера метрик: %v", err)
		}
	}()

	mongoClient := config.ConnectMongo(mongoURI)
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Ошибка отключения MongoDB: %v", err)
		}
	}()
	db := mongoClient.Database(dbName)
	if err := migrations.RunAll(db); err != nil {
		log.Fatalf("❌ Ошибка при применении миграций: %v", err)
	}
	orderRepo := repository.NewOrderRepository(db)

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Printf("Ошибка закрытия Redis: %v", err)
		}
	}()

	eventPublisher := events.NewEventPublish("amqp://guest:guest@localhost:5672/")

	userConn, err := grpc.Dial(userServiceAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не удалось подключиться к UserService: %v", err)
	}
	defer userConn.Close()
	userClient := usersproto.NewUserServiceClient(userConn)

	orderService := service.NewOrderService(orderRepo, userClient, eventPublisher, rdb)

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
