package main

import (
	"context"
	"log"
	"net"
	"net/http" // 📌 добавлено для метрик

	"github.com/prometheus/client_golang/prometheus/promhttp" // 📌 добавлено для метрик
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"gotune/events"
	"gotune/users/intern/config"
	"gotune/users/intern/repository"
	"gotune/users/intern/service"
	"gotune/users/metrics" // 📌 добавлено для метрик
	"gotune/users/migrations"
	"gotune/users/pkg/mailer"
	"gotune/users/proto"
)

const (
	mongoURI = "mongodb://localhost:27017"
	dbName   = "gotune_users"
)

func main() {
	// 📌 Запуск HTTP-сервера для метрик
	go func() {
		metrics.Register() // Регистрация метрик
		http.Handle("/metrics", promhttp.Handler())
		log.Println("📊 Метрики доступны на :2112/metrics")
		if err := http.ListenAndServe(":2112", nil); err != nil {
			log.Fatalf("Ошибка запуска HTTP-сервера метрик: %v", err)
		}
	}()

	mongoClient := config.ConnectMon(mongoURI)
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Ошибка отключения MongoDB: %v", err)
		}
	}()
	db := mongoClient.Database(dbName)

	if err := migrations.RunAll(db); err != nil {
		log.Fatalf("❌ Ошибка при применении миграций: %v", err)
	}

	userRepo := repository.NewUserRepository(db)

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Printf("Ошибка закрытия Redis: %v", err)
		}
	}()

	eventPublisher := events.NewEventPublish("amqp://guest:guest@localhost:5672/")

	emailSender := &mailer.SMTPMailer{
		From:     "berikbakhtiarovich@gmail.com",
		Host:     "smtp.gmail.com",
		Port:     587,
		Username: "berikbakhtiarovich@gmail.com",
		Password: "rtbdxwkqjdwrowsm",
	}

	userService := service.NewUserService(userRepo, eventPublisher, rdb, emailSender)

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
