package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type ApiResponseChat struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Message   struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}

func mainChat() {
	body := []byte(`{
        "model": "llama2",
        "messages": [
            { "role": "user", "content": "como redirecciono una ruta con nginx" }
        ]
    }`)
	resp, err := http.Post("https://ollama.migpt.cl/api/chat", "application/json", bytes.NewBuffer(body))

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	var fullResponse string

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var response ApiResponseChat
		if err := json.Unmarshal(scanner.Bytes(), &response); err != nil {
			log.Fatal(err)
		}

		fullResponse += response.Message.Content

		if response.Done {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(fullResponse)
}
