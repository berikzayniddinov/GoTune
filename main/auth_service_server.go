// main.go

package main

import (
	"log"
	"net"

	"github.com/berik/GoTune/Infrastructure/Persistence/Context"
	"github.com/berik/GoTune/Internal/Repositories"
	"github.com/berik/GoTune/Internal/Services"
	"github.com/berik/GoTune/Presentation/Controllers"
	pb "github.com/berik/GoTune/proto"
	"google.golang.org/grpc"
)

const port = ":50051"

func main() {
	// Инициализация БД
	db, err := Context.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err = Context.InitUserTable(db); err != nil {
		log.Fatalf("failed to init user table: %v", err)
	}

	userRepo := Repositories.NewPostgresUserRepository(db)
	authService := Services.NewAuthService(userRepo)

	// Создание контроллера
	authController := Controllers.NewAuthController(authService)

	// Настройка gRPC сервера
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, authController)

	// Запуск сервера
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Server started on port %s", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
