package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"quizzy/internal/models"
	"quizzy/internal/storage"
	"strconv"
	"time"
)

type RoomHandler struct {
	storage storage.Storage
}

func NewRoomHandler(storage storage.Storage) *RoomHandler {
	return &RoomHandler{storage: storage}
}

func (rh *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var room models.Room
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	_, err := rh.storage.GetQuiz(room.QuizID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	roomCode := rand.Intn(9000) + 1000
	roomCodeStr := strconv.Itoa(roomCode)
	room.RoomID = roomCodeStr
	room.CurrentQuestion = -1
	room.Status = "waiting"
	err = rh.storage.CreateRoom(room)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(room)
}

func (rh *RoomHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	var joinRequest models.JoinRequest
	if err := json.NewDecoder(r.Body).Decode(&joinRequest); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	roomID := r.PathValue("code")
	room, err := rh.storage.GetRoom(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if room.Status != "waiting" {
		http.Error(w, "entry closed", http.StatusConflict)
		return
	}
	playerID := time.Now().Unix()
	playerIDStr := strconv.FormatInt(playerID, 10)
	player := models.Player{
		ID:    playerIDStr,
		Name:  joinRequest.Name,
		Score: 0,
	}
	if err = rh.storage.AddPlayer(roomID, player); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(player)
}

func (rh *RoomHandler) NextQuestion(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("code")
	room, err := rh.storage.GetRoom(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if room.Status == "finished" {
		http.Error(w, "game over", http.StatusConflict)
		return
	}
	quiz, err := rh.storage.GetQuiz(room.QuizID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	room.CurrentQuestion++
	var publicQuestion *models.PublicQuestion
	if room.CurrentQuestion >= len(quiz.Questions) {
		room.Status = "finished"
	} else {
		room.Status = "active"
		room.StartTime = time.Now().Unix()
		publicQuestion = &models.PublicQuestion{
			Text:    quiz.Questions[room.CurrentQuestion].Text,
			Options: quiz.Questions[room.CurrentQuestion].Options,
		}
	}
	if err = rh.storage.UpdateRoom(room); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	nextQuestion := models.NextQuestionResponse{
		Status:          room.Status,
		CurrentQuestion: room.CurrentQuestion,
		Question:        publicQuestion,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nextQuestion)
}
