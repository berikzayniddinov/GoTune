package client

import (
	"google.golang.org/grpc"
	"gotune/users/proto"
	"log"
)

func NewUserServiceClient(connStr string) proto.UserServiceClient {
	conn, err := grpc.Dial(connStr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не удалось подключиться к UsersService: %v", err)
	}
	return proto.NewUserServiceClient(conn)
}
