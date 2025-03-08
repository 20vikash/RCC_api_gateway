package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("read error:", err)
			break
		}
		fmt.Printf("Received: %s\n", msg)

		if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
			fmt.Println("write error:", err)
			break
		}
	}
}
