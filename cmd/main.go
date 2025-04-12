package main

import (
	"fmt"
	"net/http"
	"room/grpc/client/ai"
	"room/internal/env"
	"room/internal/mq"

	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Application struct {
	Port      string
	AIService ai.AIServiceClient
	MqChannel *amqp.Channel
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	mq := &mq.MQ{
		User: env.GetMqUser(),
		Pass: env.GetMqPassword(),
		Port: "5672",
	}

	con := mq.ConnectToMq()

	app := &Application{
		Port:      "6969",
		AIService: ai.ConnectToAIService(),
		MqChannel: mq.CreateChannel(con),
	}

	fs := http.FileServer(http.Dir("../web"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/generate", app.generate)
	http.HandleFunc("/debug", app.debugCode)
	http.HandleFunc("/output", app.outputCode)
	http.HandleFunc("/createroom", createRoom)
	http.HandleFunc("/join", joinRoom)
	http.HandleFunc("/result", app.ResultQ)
	http.HandleFunc("/test", callbackHandler)

	fmt.Println("Server started on http://localhost:6969")
	err := http.ListenAndServe(":"+app.Port, nil)
	if err != nil {
		fmt.Println("ListenAndServe Error:", err)
	}
}
