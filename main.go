package main

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/genai"
	"log"
	"net/http"
	"os"
)

type Request struct {
	Text string `json:"text"`
}

type Response struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

var (
	client *genai.Client
	ctx    context.Context
)

func initGeminiClient() error {
	apiKey := "AIzaSyD8CLjZ597PM4bwEgyRzNIdqmM6FzyEEbw"

	ctx = context.Background()
	var err error
	client, err = genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	return err
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := Response{}
	result, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash", genai.Text(req.Text), nil)
	if err != nil {
		resp.Error = err.Error()
	} else {
		resp.Result = fmt.Sprint(result)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	if err := initGeminiClient(); err != nil {
		log.Fatalf("Failed to initialize Gemini client: %v", err)
	}

	http.HandleFunc("/api/generate", handleRequest)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
