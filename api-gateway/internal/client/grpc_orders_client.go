package client

import (
	"google.golang.org/grpc"
	"gotune/order/proto"
	"log"
)

func NewOrderServiceClient(connStr string) proto.OrderServiceClient {
	conn, err := grpc.Dial(connStr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не удалось подключиться к OrderService: %v", err)
	}
	return proto.NewOrderServiceClient(conn)
}
