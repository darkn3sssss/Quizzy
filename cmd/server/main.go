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
	roomHandler := handlers.NewRoomHandler(memory)
	answerHandler := handlers.NewAnswerHandler(memory)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HealthStatus{Status: "ok"})
	})
	mux.HandleFunc("GET /quiz/{id}", quizHandler.GetQuiz)
	mux.HandleFunc("GET /rooms/{code}/state", roomHandler.GetState)
	mux.HandleFunc("POST /quiz", quizHandler.PostQuiz)
	mux.HandleFunc("POST /rooms", roomHandler.CreateRoom)
	mux.HandleFunc("POST /rooms/{code}/join", roomHandler.JoinRoom)
	mux.HandleFunc("POST /rooms/{code}/next", roomHandler.NextQuestion)
	mux.HandleFunc("POST /rooms/{code}/answer", answerHandler.AnswerOnQusetion)
	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/", fileServer)
	log.Println("сервер слушает порт :8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(mux)))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
