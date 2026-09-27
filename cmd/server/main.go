package main

import (
	"encoding/json"
	"log"
	"net/http"
	"quizzy/internal/handlers"
	"quizzy/internal/storage"
)

type HealthStatus struct {
	Status string `json:"status"`
}

func main() {
	memory := storage.InitMemStorage()
	quizHandler := handlers.NewQuizHandler(memory)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HealthStatus{Status: "ok"})
	})
	mux.HandleFunc("GET /quiz/{id}", quizHandler.GetQuiz)
	mux.HandleFunc("POST /quiz", quizHandler.PostQuiz)
	log.Println("сервер слушает порт :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
