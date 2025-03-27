package output

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectToAIService() *grpc.ClientConn {
	conn, err := grpc.NewClient("localhost:6971", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}

	return conn
}
