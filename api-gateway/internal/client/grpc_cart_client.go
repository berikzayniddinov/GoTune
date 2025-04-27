package client

import (
	"google.golang.org/grpc"
	"gotune/cart/proto"
	"log"
)

func NewCartServiceClient(connStr string) proto.CartServiceClient {
	conn, err := grpc.Dial(connStr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не удалось подключиться к CartService: %v", err)
	}
	return proto.NewCartServiceClient(conn)
}
