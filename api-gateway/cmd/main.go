package main

import (
	"gotune/api-gateway/internal/client"
	"gotune/api-gateway/internal/handler"
	"log"
	"net/http"
)

func main() {
	userClient := client.NewUserServiceClient("localhost:50051")
	userHandler := handler.NewUserHandler(userClient)

	http.HandleFunc("/register", userHandler.Register)
	http.HandleFunc("/login", userHandler.Login)
	http.HandleFunc("/users", userHandler.GetAllUsers)

	log.Println("API Gateway запущен на порту 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Ошибка запуска API Gateway: %v", err)
	}
}
