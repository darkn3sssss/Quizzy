package storage

import (
	"errors"
	"quizzy/internal/models"
	"sync"
)

type MemoryStorage struct {
	quizzes map[string]*models.Quiz
	rooms   map[string]*models.Room
	Mu      *sync.RWMutex
}

func InitMemStorage() MemoryStorage {
	return MemoryStorage{
		quizzes: make(map[string]*models.Quiz),
		rooms:   make(map[string]*models.Room),
	}
}

func (memS MemoryStorage) CreateQuiz(q models.Quiz) error {
	memS.Mu.RLock()
	if _, ok := memS.quizzes[q.ID]; ok {
		return errors.New("quiz with this id exist")
	}
	memS.Mu.RUnlock()
	memS.Mu.Lock()
	memS.quizzes[q.ID] = &q
	memS.Mu.Unlock()
	return nil
}

func (memS MemoryStorage) GetQuiz(id string) (models.Quiz, error) {
	memS.Mu.RLock()
	if _, ok := memS.quizzes[id]; !ok {
		return models.Quiz{}, errors.New("invalid quiz ID")
	}
	memS.Mu.RUnlock()
	return *memS.quizzes[id], nil
}

func (memS MemoryStorage) ListQuizzes() ([]models.Quiz, error) {
	listQuizzes := []models.Quiz{}
	memS.Mu.RLock()
	for _, quiz := range memS.quizzes {
		listQuizzes = append(listQuizzes, *quiz)
	}
	memS.Mu.RUnlock()
	return listQuizzes, nil
}

func (memS MemoryStorage) CreateRoom(r models.Room) error {
	memS.Mu.RLock()
	if _, ok := memS.rooms[r.RoomID]; ok {
		return errors.New("room with this id exist")
	}
	memS.Mu.RUnlock()
	memS.Mu.Lock()
	memS.rooms[r.RoomID] = &r
	memS.rooms[r.RoomID].Answers = make(map[string][]models.Answer)
	memS.rooms[r.RoomID].Players = make(map[string]*models.Player)
	memS.Mu.Unlock()
	return nil
}

func (memS MemoryStorage) GetRoom(id string) (models.Room, error) {
	memS.Mu.RLock()
	if _, ok := memS.rooms[id]; !ok {
		return models.Room{}, errors.New("invalid room ID")
	}
	memS.Mu.RUnlock()
	return *memS.rooms[id], nil
}

func (memS MemoryStorage) UpdateRoom(r models.Room) error {
	memS.Mu.RLock()
	if _, ok := memS.rooms[r.RoomID]; !ok {
		return errors.New("cannot update a non-existent room")
	}
	memS.Mu.RUnlock()
	memS.Mu.Lock()
	memS.rooms[r.RoomID] = &r
	memS.Mu.Unlock()
	return nil
}

func (memS MemoryStorage) AddPlayer(roomID string, p models.Player) error {
	memS.Mu.RLock()
	if _, ok := memS.rooms[roomID]; !ok {
		return errors.New("invalid room ID")
	}
	memS.Mu.RUnlock()
	memS.Mu.Lock()
	memS.rooms[roomID].Players[p.ID] = &p
	memS.Mu.Unlock()
	return nil
}

func (memS MemoryStorage) SaveAnswer(roomID string, a models.Answer) error {
	memS.Mu.RLock()
	if _, ok := memS.rooms[roomID]; !ok {
		return errors.New("invalid room ID")
	}
	memS.Mu.RUnlock()
	memS.Mu.Lock()
	memS.rooms[roomID].Answers[a.QuestionID] = append(memS.rooms[roomID].Answers[a.QuestionID], a)
	memS.Mu.Unlock()
	return nil
}
