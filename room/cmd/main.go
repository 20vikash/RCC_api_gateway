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
		Port: "6969",
	}

	fs := http.FileServer(http.Dir("../web"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/generate", generate)
	http.HandleFunc("/debug", debugCode)
	http.HandleFunc("/output", output)

	fmt.Println("Server started on http://localhost:6969")
	err := http.ListenAndServe(":"+app.Port, nil)
	if err != nil {
		fmt.Println("ListenAndServe Error:", err)
	}
}
