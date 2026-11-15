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
	room.Status = "waiting"
	room.CountFinishPlayer = 0
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
	playerID := time.Now().UnixNano()
	playerIDStr := strconv.FormatInt(playerID, 10)
	player := models.Player{
		ID:              playerIDStr,
		Name:            joinRequest.Name,
		Status:          true,
		Score:           0,
		CurrentQuestion: -1,
		StartTime:       0,
	}
	_, err := rh.storage.UpdateRoom(roomID, func(room *models.Room) error {
		if room.Status != "waiting" {
			return roomError{http.StatusConflict, "entry closed"}
		}
		room.Players[player.ID] = &player
		return nil
	})
	if err != nil {
		writeRoomError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

func (rh *RoomHandler) StartQuizz(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("code")
	room, err := rh.storage.UpdateRoom(roomID, func(room *models.Room) error {
		if room.Status == "waiting" {
			room.Status = "active"
		}
		return nil
	})
	if err != nil {
		writeRoomError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
}

func (rh *RoomHandler) NextQuestionForPlayer(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("code")
	room, err := rh.storage.GetRoom(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	quiz, err := rh.storage.GetQuiz(room.QuizID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var nextRequest models.NextQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&nextRequest); err != nil || nextRequest.QuestionIndex == nil {
		http.Error(w, "question_index is required", http.StatusBadRequest)
		return
	}
	var nextQuestion models.NextQuestionResponse
	_, err = rh.storage.UpdateRoom(roomID, func(room *models.Room) error {
		if room.Status != "active" {
			return roomError{http.StatusConflict, "player cannot answer"}
		}
		id := r.PathValue("id")
		if _, ok := room.Players[id]; !ok {
			return roomError{http.StatusNotFound, "player not found"}
		}
		player := room.Players[id]
		if player.Status == false {
			return roomError{http.StatusConflict, "player already finished"}
		}
		if *nextRequest.QuestionIndex != player.CurrentQuestion {
			return roomError{http.StatusConflict, "stale question"}
		}
		if player.CurrentQuestion < -1 || player.CurrentQuestion >= len(quiz.Questions) {
			return roomError{http.StatusConflict, "CurrentQuestion out of range"}
		}
		if player.CurrentQuestion >= 0 {
			question := quiz.Questions[player.CurrentQuestion]
			answered := false
			for _, answer := range room.Answers[question.ID] {
				if answer.PlayerID == id {
					answered = true
					break
				}
			}
			if !answered && time.Now().Unix()-player.StartTime < int64(question.TimeLimit) {
				return roomError{http.StatusConflict, "answer the current question first"}
			}
		}
		player.CurrentQuestion++
		var publicQuestion *models.PublicQuestion
		if player.CurrentQuestion >= len(quiz.Questions) {
			player.Status = false
			room.CountFinishPlayer++
			if room.CountFinishPlayer == len(room.Players) {
				room.Status = "finished"
			}
		} else {
			player.StartTime = time.Now().Unix()
			publicQuestion = &models.PublicQuestion{
				Text:    quiz.Questions[player.CurrentQuestion].Text,
				Options: quiz.Questions[player.CurrentQuestion].Options,
			}
		}
		nextQuestion = models.NextQuestionResponse{
			CurrentQuestion: player.CurrentQuestion,
			Question:        publicQuestion,
		}
		return nil
	})
	if err != nil {
		writeRoomError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nextQuestion)
}

func (rh *RoomHandler) GetStateForPlayer(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("code")
	room, err := rh.storage.GetRoom(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	quiz, err := rh.storage.GetQuiz(room.QuizID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id := r.PathValue("id")
	if _, ok := room.Players[id]; !ok {
		http.Error(w, "player not found", http.StatusNotFound)
		return
	}
	player := room.Players[id]
	elapsed := time.Now().Unix() - player.StartTime
	var timeRemaining int64
	var publicQuestion *models.PublicQuestion
	if -1 < player.CurrentQuestion && player.CurrentQuestion < len(quiz.Questions) {
		publicQuestion = &models.PublicQuestion{
			Text:    quiz.Questions[player.CurrentQuestion].Text,
			Options: quiz.Questions[player.CurrentQuestion].Options,
		}
		timeRemaining = int64(quiz.Questions[player.CurrentQuestion].TimeLimit) - elapsed
	} else {
		timeRemaining = 0
	}
	if timeRemaining < 0 {
		timeRemaining = 0
	}
	roomStateResponse := models.RoomPlayerStateResponse{
		Question:        publicQuestion,
		CurrentQuestion: player.CurrentQuestion,
		TimeRemaining:   timeRemaining,
		Player:          *player,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roomStateResponse)
}

func (rh *RoomHandler) GetState(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("code")
	room, err := rh.storage.GetRoom(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	roomStateResponse := models.RoomStateResponse{
		Status:  room.Status,
		Players: room.Players,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roomStateResponse)
}

// func (rh *RoomHandler) NextQuestion(w http.ResponseWriter, r *http.Request) {
// 	roomID := r.PathValue("code")
// 	room, err := rh.storage.GetRoom(roomID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}
// 	if room.Status == "finished" {
// 		http.Error(w, "game over", http.StatusConflict)
// 		return
// 	}
// 	quiz, err := rh.storage.GetQuiz(room.QuizID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	room.CurrentQuestion++
// 	var publicQuestion *models.PublicQuestion
// 	if room.CurrentQuestion >= len(quiz.Questions) {
// 		room.Status = "finished"
// 	} else {
// 		room.Status = "active"
// 		room.StartTime = time.Now().Unix()
// 		publicQuestion = &models.PublicQuestion{
// 			Text:    quiz.Questions[room.CurrentQuestion].Text,
// 			Options: quiz.Questions[room.CurrentQuestion].Options,
// 		}
// 	}
// 	if err = rh.storage.UpdateRoom(room); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	nextQuestion := models.NextQuestionResponse{
// 		Status:          room.Status,
// 		CurrentQuestion: room.CurrentQuestion,
// 		Question:        publicQuestion,
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(nextQuestion)
// }
