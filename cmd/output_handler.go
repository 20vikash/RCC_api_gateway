package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	o "room/grpc/client/output"
)

type OutputResponse struct {
	Output string
}

func (app *Application) outputCode(w http.ResponseWriter, r *http.Request) {
	language := r.URL.Query().Get("language")
	roomId := r.URL.Query().Get("id")
	userName := r.URL.Query().Get("username")

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	var codeData Code
	var output OutputResponse

	var res *o.OutputResponse
	var err error

	if err = json.NewDecoder(r.Body).Decode(&codeData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	code := codeData.Code

	req := &o.OutputRequest{
		RoomID:   roomId,
		UserName: userName,
		Language: language,
		Code:     code,
	}

	if language == "cpp" || language == "c" {
		res, err = app.OutputService.OutputCCpp(ctx, req)
		if err != nil {
			log.Println(err)
		}
	} else if language == "python" {
		res, err = app.OutputService.OutputPython(ctx, req)
		if err != nil {
			log.Println(err)
		}
	} else if language == "go" || language == "php" {
		res, err = app.OutputService.OutputGolangPHP(ctx, req)
		if err != nil {
			log.Println(err)
		}
	}

	output = OutputResponse{Output: res.Message}
	j, err := json.Marshal(output)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusInternalServerError)
	}

	w.Write(j)
}
