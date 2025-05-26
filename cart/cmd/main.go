package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"

	"gotune/cart/internal/config"
	"gotune/cart/internal/repository"
	"gotune/cart/internal/service"
	"gotune/cart/metrics"
	"gotune/cart/migrations"
	"gotune/cart/proto"
	"gotune/events"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	mongoURI = "mongodb://localhost:27017"
	dbName   = "gotune_cart"
)

func main() {
	// 📌 Запуск HTTP-сервера для метрик
	go func() {
		metrics.Register() // Регистрация метрик
		http.Handle("/metrics", promhttp.Handler())
		log.Println("📊 Метрики доступны на :2115/metrics")
		if err := http.ListenAndServe(":2115", nil); err != nil {
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
	if err := migrations.RunM(db); err != nil {
		log.Fatalf("❌ Ошибка при применении миграций: %v", err)
	}
	cartRepo := repository.NewCartRepositories(db)

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Printf("Ошибка закрытия Redis: %v", err)
		}
	}()

	eventPublisher := events.NewEventPublish("amqp://guest:guest@localhost:5672/")

	cartService := service.NewCartService(cartRepo, eventPublisher, rdb)

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
