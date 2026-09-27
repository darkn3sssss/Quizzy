package handlers

import (
	"encoding/json"
	"net/http"
	"quizzy/internal/models"
	"quizzy/internal/storage"
	"time"
)

type AnswerHandler struct {
	storage storage.Storage
}

func NewAnswerHandler(storage storage.Storage) *AnswerHandler {
	return &AnswerHandler{storage: storage}
}

func (ah *AnswerHandler) AnswerOnQusetion(w http.ResponseWriter, r *http.Request) {
	var answerRequest models.AnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&answerRequest); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	roomID := r.PathValue("code")
	room, err := ah.storage.GetRoom(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if room.Status != "active" {
		http.Error(w, "room not active", http.StatusConflict)
		return
	}
	if _, ok := room.Players[answerRequest.PlayerID]; !ok {
		http.Error(w, "player not found", http.StatusNotFound)
		return
	}
	quiz, err := ah.storage.GetQuiz(room.QuizID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	currentIDQuestion := quiz.Questions[room.CurrentQuestion].ID
	for _, answer := range room.Answers[currentIDQuestion] {
		if answer.PlayerID == answerRequest.PlayerID {
			http.Error(w, "player answered", http.StatusConflict)
			return
		}
	}
	answeredAt := time.Now().Unix()
	newAnswer := models.Answer{
		PlayerID:   answerRequest.PlayerID,
		QuestionID: currentIDQuestion,
		Option:     answerRequest.Option,
		AnsweredAt: answeredAt,
	}
	if err = ah.storage.SaveAnswer(room.RoomID, newAnswer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	correctOption := quiz.Questions[room.CurrentQuestion].CorrectOption
	isCorrect := false
	if answerRequest.Option == correctOption {
		room.Players[answerRequest.PlayerID].Score += 100
		isCorrect = true
	}
	if err = ah.storage.UpdateRoom(room); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newScore := room.Players[answerRequest.PlayerID].Score
	answerResponse := models.AnswerResponse{
		Correct:  isCorrect,
		NewScore: newScore,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answerResponse)
}
