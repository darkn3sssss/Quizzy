package handlers

import (
	"encoding/json"
	"net/http"
	"quizzy/internal/models"
	"quizzy/internal/storage"
	"strconv"
	"time"
)

type QuizHandler struct {
	storage storage.Storage
}

func NewQuizHandler(storage storage.Storage) *QuizHandler {
	return &QuizHandler{storage: storage}
}

func (qh *QuizHandler) GetQuiz(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	quiz, err := qh.storage.GetQuiz(strID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quiz)
}

func (qh *QuizHandler) ListQuizzes(w http.ResponseWriter, r *http.Request) {
	listQuizzes, _ := qh.storage.ListQuizzes()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listQuizzes)
}

func (qh *QuizHandler) PostQuiz(w http.ResponseWriter, r *http.Request) {
	var quiz models.Quiz
	if err := json.NewDecoder(r.Body).Decode(&quiz); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	unixTime := time.Now().UnixNano()
	quizID := strconv.FormatInt(unixTime, 10)
	quiz.ID = quizID
	if err := qh.storage.CreateQuiz(quiz); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(quiz)
}
