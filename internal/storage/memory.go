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
	quizzes map[string]*models.Quiz
	rooms   map[string]*models.Room
}

func InitMemStorage() MemoryStorage {
	return MemoryStorage{
		quizzes: make(map[string]*models.Quiz),
		rooms:   make(map[string]*models.Room),
	}
}

func (memS MemoryStorage) CreateQuiz(q models.Quiz) error {
	if _, ok := memS.quizzes[q.ID]; ok {
		return errors.New("quiz with this id exist")
	}
	memS.quizzes[q.ID] = &q
	return nil
}

func (memS MemoryStorage) GetQuiz(id string) (models.Quiz, error) {
	if _, ok := memS.quizzes[id]; !ok {
		return models.Quiz{}, errors.New("wrong quiz ID")
	}
	return *memS.quizzes[id], nil
}

func (memS MemoryStorage) CreateRoom(r models.Room) error {
	if _, ok := memS.rooms[r.RoomCode]; ok {
		return errors.New("room with this id exist")
	}
	memS.rooms[r.RoomCode] = &r
	memS.rooms[r.RoomCode].Answers = make(map[string][]models.Answer)
	return nil
}

func (memS MemoryStorage) GetRoom(id string) (models.Room, error) {
	if _, ok := memS.rooms[id]; !ok {
		return models.Room{}, errors.New("wrong room ID")
	}
	return *memS.rooms[id], nil
}

func (memS MemoryStorage) UpdateRoom(r models.Room) error {
	if _, ok := memS.rooms[r.RoomCode]; !ok {
		return errors.New("cannot update a non-existent room")
	}
	memS.rooms[r.RoomCode] = &r
	return nil
}

func (memS MemoryStorage) AddPlayer(roomID string, p models.Player) error {
	if _, ok := memS.rooms[roomID]; !ok {
		return errors.New("wrong room ID")
	}
	memS.rooms[roomID].Players = append(memS.rooms[roomID].Players, p)
	return nil
}

func (memS MemoryStorage) SaveAnswer(roomID string, a models.Answer) error {
	if _, ok := memS.rooms[roomID]; !ok {
		return errors.New("wrong room ID")
	}
	memS.rooms[roomID].Answers[a.QuestionID] = append(memS.rooms[roomID].Answers[a.QuestionID], a)
	return nil
}
