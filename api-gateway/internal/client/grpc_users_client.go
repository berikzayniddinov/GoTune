package client

import (
	"google.golang.org/grpc"
	"gotune/instruments/proto"
	"log"
)

func NewInstrumentServiceClient(connStr string) proto.InstrumentServiceClient {
	conn, err := grpc.Dial(connStr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не удалось подключиться к InstrumentService: %v", err)
	}
	return proto.NewInstrumentServiceClient(conn)
}
