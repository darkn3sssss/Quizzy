package storage

import (
	"errors"
	"quizzy/internal/models"
)

type Storage interface {
	CreateQuiz(q models.Quiz) error
	GetQuiz(id string) (models.Quiz, error)
	CreateRoom(r models.Room) error
	GetRoom(id string) (models.Room, error)
	UpdateRoom(r models.Room) error
	AddPlayer(roomID string, p models.Player) error
	SaveAnswer(roomID string, a models.Answer) error
}

type MemoryStorage struct {
	quizzes map[string]models.Quiz
	rooms   map[string]models.Room
}

func InitMemStorage() MemoryStorage {
	return MemoryStorage{
		quizzes: make(map[string]models.Quiz),
		rooms:   make(map[string]models.Room),
	}
}

func (memS MemoryStorage) CreateQuiz(q models.Quiz) error {
	if _, ok := memS.quizzes[q.ID]; ok {
		return errors.New("Quiz with this id exist")
	}
	memS.quizzes[q.ID] = q
	return nil
}

func (memS MemoryStorage) GetQuiz(id string) (models.Quiz, error) {
	if _, ok := memS.quizzes[id]; !ok {
		return models.Quiz{}, errors.New("Wrong quiz ID")
	}
	return memS.quizzes[id], nil
}

func (memS MemoryStorage) CreateRoom(r models.Room) error {
	if _, ok := memS.rooms[r.QuizID]; ok {
		return errors.New("Room with this id exist")
	}
	memS.rooms[r.QuizID] = r
	return nil
}

func (memS MemoryStorage) GetRoom(id string) (models.Room, error) {
	if _, ok := memS.rooms[id]; !ok {
		return models.Room{}, errors.New("Wrong room ID")
	}
	return memS.rooms[id], nil
}

func (memS MemoryStorage) UpdateRoom(r models.Room) error {
	return nil
}

func (memS MemoryStorage) AddPlayer(roomID string, p models.Player) error {
	if _, ok := memS.rooms[roomID]; !ok {
		return errors.New("Wrong room ID")
	}
	ptrPlayers := memS.rooms[roomID].Players
	ptrPlayers = append(ptrPlayers, p)
	return nil
}

func (memS MemoryStorage) SaveAnswer(roomID string, a models.Answer) error {
	if _, ok := memS.rooms[roomID]; !ok {
		return errors.New("Wrong room ID")
	}
	ptrAnswers := memS.rooms[roomID].Answers
	ptrAnswers[a.QuestionID] = append(ptrAnswers[a.QuestionID], a)
	return nil
}
