package ai

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectToAIService() AIServiceClient {
	conn, err := grpc.NewClient("viksync_ai:6970", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}

	cl := NewAIServiceClient(conn)

	return cl
}
