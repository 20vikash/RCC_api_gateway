package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Application struct {
	Port string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	app := &Application{
		Port: ":69",
	}

	http.HandleFunc("/ws", handleConnections)

	fmt.Println("WebSocket server started on :69")
	err := http.ListenAndServe(app.Port, nil)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
