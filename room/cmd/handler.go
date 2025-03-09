package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var roomConns = make(map[string][]*websocket.Conn)
var changes = make(map[string][]string)

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
			fmt.Println("read error:", err)
			break
		}
		fmt.Printf("Received: %s\n", msg)
		m := string(msg)

		s := strings.Split(m, "~")
		room, message := s[0], s[1]

		if message == "load" {
			if len(changes[room]) > 0 {
				history, _ := json.Marshal(changes[room])
				if err := ws.WriteMessage(websocket.TextMessage, history); err != nil {
					fmt.Println("write error:", err)
					return
				}
				continue
			}
		}

		rcns := roomConns[room]

		if message != "load" {
			changes[roomID] = append(changes[roomID], message)
		}

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

type CodeResponse struct {
	Code string `json:"code"`
}

func generate(w http.ResponseWriter, r *http.Request) {
	prompt := r.URL.Query().Get("prompt")
	langauge := r.URL.Query().Get("language")
	roomID := r.URL.Query().Get("id")

	changes[roomID] = changes[roomID][:0]

	cmd := exec.Command("python3", "../codellama/generate.py", prompt, langauge)

	var stdOut bytes.Buffer
	cmd.Stdout = &stdOut

	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}

	res := &CodeResponse{
		Code: stdOut.String(),
	}

	jm, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
	}

	w.Write(jm)
}

func debugCode(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("id")
	language := r.URL.Query().Get("language")

	var codeData CodeResponse

	changes[roomID] = changes[roomID][:0]

	if err := json.NewDecoder(r.Body).Decode(&codeData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	cmd := exec.Command("python3", "../codellama/debug.py", codeData.Code, language)

	var stdOut bytes.Buffer
	cmd.Stdout = &stdOut

	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}

	res := &CodeResponse{
		Code: stdOut.String(),
	}

	jm, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
	}

	w.Write(jm)
}

type OutputResponse struct {
	Output string
}

func output(w http.ResponseWriter, r *http.Request) {
	language := r.URL.Query().Get("language")
	roomId := r.URL.Query().Get("id")
	userName := r.URL.Query().Get("username")

	var codeData CodeResponse
	var output OutputResponse

	var res string

	if err := json.NewDecoder(r.Body).Decode(&codeData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	code := codeData.Code

	if language == "cpp" || language == "c" {
		res = outputCCpp(roomId, userName, language, code)
	} else if language == "python" {
		res = outputPython(code)
	}

	output = OutputResponse{Output: res}
	j, err := json.Marshal(output)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusInternalServerError)
	}

	w.Write(j)
}

func outputCCpp(roomId string, userName string, language string, code string) string {
	var cmd *exec.Cmd
	var stdErr bytes.Buffer
	var stdOut bytes.Buffer
	var res string

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if language == "cpp" {
		language = "c++"
	}

	fileName := roomId + userName

	cmd = exec.Command("g++", "-x", language, "-", "-o", fileName)
	cmd.Stdin = bytes.NewBufferString(code)
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		res = stdErr.String()
		return res
	}

	defer os.Remove(fileName)

	cmd = exec.CommandContext(ctx, "./"+fileName)
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err = cmd.Run()

	if err != nil {
		if err == err.(*exec.ExitError) {
			res = "Took too long to generate the output"
			return res
		}
		log.Println(err.Error())
		res = stdErr.String()
		return res
	} else {
		res = stdOut.String()
		return res
	}
}

func outputPython(code string) string {
	var stdErr bytes.Buffer
	var stdOut bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "python3")
	cmd.Stdin = bytes.NewBufferString(code)
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()

	if err != nil {
		if err == err.(*exec.ExitError) && stdErr.String() == "" {
			return "Took too long to generate the output"
		}
		return stdErr.String()
	}

	return stdOut.String()
}
