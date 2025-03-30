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
	"time"
)

type OutputResponse struct {
	Output string
}

func (app *Application) outputCode(w http.ResponseWriter, r *http.Request) {
	language := r.URL.Query().Get("language")
	roomId := r.URL.Query().Get("id")
	userName := r.URL.Query().Get("username")

	var codeData Code
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
	} else if language == "go" || language == "php" {
		res = outputGolangPHP(roomId, userName, code, language)
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

func outputGolangPHP(roomID string, userName string, code string, language string) string {
	var stdErr bytes.Buffer
	var stdOut bytes.Buffer

	var extension string
	var cmd *exec.Cmd
	var file string
	var filePath string

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	file = roomID + userName

	if language == "go" {
		extension = ".go"
		filePath = "../sandbox/" + file + extension
		cmd = exec.CommandContext(ctx, "go", "run", filePath)
	} else {
		extension = ".php"
		filePath = "../sandbox/" + file + extension
		cmd = exec.CommandContext(ctx, "php", filePath)
	}

	if err := os.WriteFile(filePath, []byte(code), 0644); err != nil {
		fmt.Println("Error writing file:", err)
		return err.Error()
	}

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error running file:", err)

		if err := os.Remove(filePath); err != nil {
			fmt.Println("Error removing file:", err)
		}

		return stdErr.String()
	}

	if err := os.Remove(filePath); err != nil {
		fmt.Println("Error removing file:", err)
	}

	return stdOut.String()
}
