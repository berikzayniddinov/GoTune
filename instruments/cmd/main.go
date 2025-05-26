package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"

	"gotune/events"
	"gotune/instruments/metrics"
	"gotune/instruments/migrations"

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
	// 📌 Запуск HTTP-сервера для метрик
	go func() {
		metrics.Register() // Регистрация метрик
		http.Handle("/metrics", promhttp.Handler())
		log.Println("📊 Метрики доступны на :2114/metrics")
		if err := http.ListenAndServe(":2114", nil); err != nil {
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
	instrumentRepo := repository.NewInstrumentRepositories(db)
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
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
