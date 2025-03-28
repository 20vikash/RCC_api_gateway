package output

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectToOutputService() OutputServiceClient {
	conn, err := grpc.NewClient("localhost:6971", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}

	cl := NewOutputServiceClient(conn)

	return cl
}
