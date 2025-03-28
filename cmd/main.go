package main

import (
	"fmt"
	"net/http"
	"room/grpc/client/ai"
	"room/grpc/client/output"

	"github.com/gorilla/websocket"
)

type Application struct {
	Port          string
	AIService     ai.AIServiceClient
	OutputService output.OutputServiceClient
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	app := &Application{
		Port:          "6969",
		AIService:     ai.ConnectToAIService(),
		OutputService: output.ConnectToOutputService(),
	}

	fs := http.FileServer(http.Dir("../web"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/generate", app.generate)
	http.HandleFunc("/debug", app.debugCode)
	http.HandleFunc("/output", outputCode)
	http.HandleFunc("/createroom", createRoom)
	http.HandleFunc("/join", joinRoom)

	fmt.Println("Server started on http://localhost:6969")
	err := http.ListenAndServe(":"+app.Port, nil)
	if err != nil {
		fmt.Println("ListenAndServe Error:", err)
	}
}
