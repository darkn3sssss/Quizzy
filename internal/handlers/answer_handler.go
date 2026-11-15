package handlers

import (
	"encoding/json"
	"math"
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

const defaultOdd float64 = 10.0

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
	quiz, err := ah.storage.GetQuiz(room.QuizID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if answerRequest.QuestionIndex == nil {
		http.Error(w, "question_index is required", http.StatusBadRequest)
		return
	}
	var answerResponse models.AnswerResponse
	_, err = ah.storage.UpdateRoom(roomID, func(room *models.Room) error {
		if room.Status != "active" {
			return roomError{http.StatusConflict, "room not active"}
		}
		if _, ok := room.Players[answerRequest.PlayerID]; !ok {
			return roomError{http.StatusNotFound, "player not found"}
		}
		player := room.Players[answerRequest.PlayerID]
		if player.Status == false {
			return roomError{http.StatusConflict, "player already finished"}
		}
		if player.CurrentQuestion < 0 || player.CurrentQuestion >= len(quiz.Questions) {
			return roomError{http.StatusConflict, "CurrentQuestion out of range"}
		}
		if *answerRequest.QuestionIndex != player.CurrentQuestion {
			return roomError{http.StatusConflict, "stale question"}
		}
		currentIDQuestion := quiz.Questions[player.CurrentQuestion].ID
		for _, answer := range room.Answers[currentIDQuestion] {
			if answer.PlayerID == answerRequest.PlayerID {
				return roomError{http.StatusConflict, "player answered"}
			}
		}
		answeredAt := time.Now().Unix()
		elapsed := answeredAt - player.StartTime
		timeLimit := int64(quiz.Questions[player.CurrentQuestion].TimeLimit)
		if elapsed > timeLimit {
			return roomError{http.StatusConflict, "time's up"}
		}
		newAnswer := models.Answer{
			PlayerID:   answerRequest.PlayerID,
			QuestionID: currentIDQuestion,
			Option:     answerRequest.Option,
			AnsweredAt: answeredAt,
		}
		room.Answers[currentIDQuestion] = append(room.Answers[currentIDQuestion], newAnswer)
		correctOption := quiz.Questions[player.CurrentQuestion].CorrectOption
		isCorrect := false
		if answerRequest.Option == correctOption {
			odd := int(math.Round((float64(timeLimit-elapsed) * defaultOdd) / float64(timeLimit)))
			if player.Name == "Cheater" || player.Name == "Admin" || player.Name == "Semen" {
				player.Score += odd * 10000
			} else {
				player.Score += odd * 10
			}
			isCorrect = true
		}
		newScore := player.Score
		answerResponse = models.AnswerResponse{
			Correct:  isCorrect,
			NewScore: newScore,
		}
		return nil
	})
	if err != nil {
		writeRoomError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answerResponse)
}

// func (ah *AnswerHandler) AnswerOnQusetion(w http.ResponseWriter, r *http.Request) {
// 	var answerRequest models.AnswerRequest
// 	if err := json.NewDecoder(r.Body).Decode(&answerRequest); err != nil {
// 		http.Error(w, "invalid JSON", http.StatusBadRequest)
// 		return
// 	}
// 	roomID := r.PathValue("code")
// 	room, err := ah.storage.GetRoom(roomID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}
// 	if room.Status != "active" {
// 		http.Error(w, "room not active", http.StatusConflict)
// 		return
// 	}
// 	if _, ok := room.Players[answerRequest.PlayerID]; !ok {
// 		http.Error(w, "player not found", http.StatusNotFound)
// 		return
// 	}
// 	quiz, err := ah.storage.GetQuiz(room.QuizID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	currentIDQuestion := quiz.Questions[room.CurrentQuestion].ID
// 	for _, answer := range room.Answers[currentIDQuestion] {
// 		if answer.PlayerID == answerRequest.PlayerID {
// 			http.Error(w, "player answered", http.StatusConflict)
// 			return
// 		}
// 	}
// 	elapsed := time.Now().Unix() - room.StartTime
// 	timeLimit := int64(quiz.Questions[room.CurrentQuestion].TimeLimit)
// 	if elapsed > timeLimit {
// 		http.Error(w, "time's up", http.StatusConflict)
// 		return
// 	}
// 	answeredAt := time.Now().Unix()
// 	newAnswer := models.Answer{
// 		PlayerID:   answerRequest.PlayerID,
// 		QuestionID: currentIDQuestion,
// 		Option:     answerRequest.Option,
// 		AnsweredAt: answeredAt,
// 	}
// 	if err = ah.storage.SaveAnswer(room.RoomID, newAnswer); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	correctOption := quiz.Questions[room.CurrentQuestion].CorrectOption
// 	isCorrect := false
// 	if answerRequest.Option == correctOption {
// 		room.Players[answerRequest.PlayerID].Score += 100
// 		isCorrect = true
// 	}
// 	if err = ah.storage.UpdateRoom(room); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	newScore := room.Players[answerRequest.PlayerID].Score
// 	answerResponse := models.AnswerResponse{
// 		Correct:  isCorrect,
// 		NewScore: newScore,
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(answerResponse)
// }

//	odds
//	timeLimit = 10
//  elapsed   = x
//
