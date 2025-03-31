package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"room/grpc/client/ai"
	"time"
)

func (app *Application) generate(w http.ResponseWriter, r *http.Request) {
	prompt := r.URL.Query().Get("prompt")
	langauge := r.URL.Query().Get("language")

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	req := &ai.AIRequest{
		Prompt:   prompt,
		Language: langauge,
	}

	res, err := app.AIService.GenerateCode(ctx, req)
	if err != nil {
		log.Println(err)
	}

	strc := &Code{
		Code: res.Message,
	}

	jm, err := json.Marshal(strc)
	if err != nil {
		log.Println(err)
	}

	w.Write(jm)
}

func (app *Application) debugCode(w http.ResponseWriter, r *http.Request) {
	language := r.URL.Query().Get("language")

	var codeData Code

	if err := json.NewDecoder(r.Body).Decode(&codeData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	req := &ai.AIRequest{
		Language: language,
		Code:     codeData.Code,
	}

	res, err := app.AIService.DebugCode(ctx, req)
	if err != nil {
		log.Println(err)
	}

	srtc := &Code{
		Code: res.Message,
	}

	jm, err := json.Marshal(srtc)
	if err != nil {
		log.Println(err)
	}

	w.Write(jm)
}
