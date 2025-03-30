package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"slices"

	"github.com/gorilla/websocket"
)

var roomConns = make(map[string][]*websocket.Conn)
var activeRooms = make([]string, 0)

func removeElement(slice []*websocket.Conn, element *websocket.Conn) []*websocket.Conn {
	for i, v := range slice {
		if v == element {
			return slices.Delete(slice, i, i+1)
		}
	}
	return slice
}

func createRoom(w http.ResponseWriter, _ *http.Request) {
	uuid, err := exec.Command("uuidgen").Output()
	uuid = []byte(strings.TrimSpace(string(uuid)))

	if err != nil {
		log.Println(err)
	}

	activeRooms = append(activeRooms, string(uuid))

	w.Write(uuid)
}

func joinRoom(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")

	c := slices.Contains(activeRooms, roomID)

	if !c {
		http.Error(w, "Room not found", http.StatusBadRequest)
		return
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	roomID := r.URL.Query().Get("roomid")

	_, exists := roomConns[roomID]
	if !exists {
		roomConns[roomID] = []*websocket.Conn{ws}
	} else {
		roomConns[roomID] = append(roomConns[roomID], ws)
	}

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Println("Client closed the connection")
				roomConns[roomID] = removeElement(roomConns[roomID], ws)

				if len(roomConns[roomID]) == 0 {
					i := slices.Index(activeRooms, roomID)
					activeRooms = slices.Delete(activeRooms, i, i+1)
					delete(roomConns, roomID)
				}

				ws.Close()
				break
			} else {
				log.Println("Error reading message:", err)
			}
		}

		fmt.Printf("Received: %s\n", msg)
		m := string(msg)

		s := strings.Split(m, "~")
		room, message := s[0], s[1]

		if s[0] == "change" {
			newLang := s[1]

			for _, con := range roomConns[roomID] {
				response := "change" + "~" + newLang

				if err := con.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
					fmt.Println(err)
					return
				}
			}
		}

		if message == "load" {
			if len(roomConns) > 0 {
				if err := roomConns[roomID][0].WriteMessage(websocket.TextMessage, []byte("lll")); err != nil {
					fmt.Println(err)
					return
				}
			}
		}

		if s[0] == "c" {
			for _, v := range roomConns[roomID] {
				if err := v.WriteMessage(websocket.TextMessage, msg); err != nil {
					fmt.Println(err)
					return
				}
			}
		}

		if room == "lll" {
			codeResponse := CodeResponse{
				Code: message,
			}

			for _, v := range roomConns[roomID] {
				if err := v.WriteMessage(websocket.TextMessage, []byte("done~"+codeResponse.Code)); err != nil {
					fmt.Println(err)
					return
				}
			}
			continue
		}

		rcns := roomConns[room]

		for _, conn := range rcns {
			if ws != conn {
				if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					fmt.Println("write error:", err)
					break
				}
			}
		}
	}
}
