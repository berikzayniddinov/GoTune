package main

import (
	"github.com/gorilla/mux"
	"gotune/api-gateway/internal/client"
	"gotune/api-gateway/internal/handler"
	"log"
	"net/http"
)

func main() {
	userClient := client.NewUserServiceClient("localhost:50051")
	userHandler := handler.NewUserHandler(userClient)

	instrumentClient := client.NewInstrumentServiceClient("localhost:50052")
	instrumentHandler := handler.NewInstrumentHandler(instrumentClient)

	cartClient := client.NewCartServiceClient("localhost:50053")
	cartHandler := handler.NewCartHandler(cartClient)

	orderClient := client.NewOrderServiceClient("localhost:50054")
	orderHandler := handler.NewOrderHandler(orderClient, userClient)

	r := mux.NewRouter()

	r.HandleFunc("/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")

	r.HandleFunc("/instruments", instrumentHandler.CreateInstrument).Methods("POST")
	r.HandleFunc("/instruments", instrumentHandler.GetAllInstruments).Methods("GET")
	r.HandleFunc("/instruments/{id}", instrumentHandler.DeleteInstrumentByID).Methods("DELETE")
	r.HandleFunc("/instruments/{id}", instrumentHandler.UpdateInstrumentByID).Methods("PATCH")

	r.HandleFunc("/cart/add", cartHandler.AddToCart).Methods("POST")
	r.HandleFunc("/cart/get", cartHandler.GetCart).Methods("GET")
	r.HandleFunc("/cart/remove", cartHandler.RemoveFromCart).Methods("POST")
	r.HandleFunc("/cart/clear", cartHandler.ClearCart).Methods("POST")

	r.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")
	r.HandleFunc("/orders", orderHandler.GetOrders).Methods("GET")
	r.HandleFunc("/orders", orderHandler.DeleteOrder).Methods("DELETE")

	log.Println("API Gateway запущен на порту 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Ошибка запуска API Gateway: %v", err)
	}
}
