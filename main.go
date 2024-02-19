package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type ApiResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

func main() {
	mainResponse()
}

func mainResponse() {
	body := []byte(`{"model":"llama2", "prompt":"dime algo de pearl jam como si fuera un twett", "system":"responde siempre en español en maximo 10 palabras", "stream":true}`)
	resp, err := http.Post("https://ollama.migpt.cl/api/generate", "application/json", bytes.NewBuffer(body))

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	var fullResponse string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var response ApiResponse
		if err := json.Unmarshal(scanner.Bytes(), &response); err != nil {
			log.Fatal(err)
		}

		// Acumular la respuesta
		fullResponse += response.Response

		// Imprimir cada fragmento de la respuesta en tiempo real
		fmt.Print(response.Response)

		if response.Done {
			fmt.Println("\n<fin>.")
			response.Response = fullResponse // Actualizar la respuesta completa
			saveResponseToFile(response)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func saveResponseToFile(response ApiResponse) {
	currentTime := time.Now()
	fileName := fmt.Sprintf("response/%d%02d%02d%02d%02d%02d%03d.json",
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		currentTime.Hour(), currentTime.Minute(), currentTime.Second(), currentTime.Nanosecond()/1000000)

	fileContent, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(fileName, fileContent, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Respuesta guardada en: %s\n", fileName)
}
