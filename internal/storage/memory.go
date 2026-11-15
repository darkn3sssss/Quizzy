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
	memS.Mu.Lock()
	defer memS.Mu.Unlock()
	if _, ok := memS.quizzes[q.ID]; ok {
		return errors.New("quiz with this id exist")
	}
	memS.quizzes[q.ID] = &q
	return nil
}

func (memS MemoryStorage) GetQuiz(id string) (models.Quiz, error) {
	memS.Mu.RLock()
	defer memS.Mu.RUnlock()
	if _, ok := memS.quizzes[id]; !ok {
		return models.Quiz{}, errors.New("invalid quiz ID")
	}
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
	memS.Mu.Lock()
	defer memS.Mu.Unlock()
	if _, ok := memS.rooms[r.RoomID]; ok {
		return errors.New("room with this id exist")
	}
	memS.rooms[r.RoomID] = &r
	memS.rooms[r.RoomID].Players = make(map[string]*models.Player)
	memS.rooms[r.RoomID].Answers = make(map[string]([]models.Answer))
	return nil
}

func (memS MemoryStorage) GetRoom(id string) (models.Room, error) {
	memS.Mu.RLock()
	defer memS.Mu.RUnlock()
	if _, ok := memS.rooms[id]; !ok {
		return models.Room{}, errors.New("invalid room ID")
	}
	return cloneRoom(*memS.rooms[id]), nil
}

func (memS MemoryStorage) UpdateRoom(id string, update func(*models.Room) error) (models.Room, error) {
	memS.Mu.Lock()
	defer memS.Mu.Unlock()
	room, ok := memS.rooms[id]
	if !ok {
		return models.Room{}, errors.New("invalid room ID")
	}
	updated := cloneRoom(*room)
	if err := update(&updated); err != nil {
		return models.Room{}, err
	}
	stored := cloneRoom(updated)
	memS.rooms[id] = &stored
	return updated, nil
}

func cloneRoom(room models.Room) models.Room {
	players := make(map[string]*models.Player, len(room.Players))
	for id, player := range room.Players {
		copy := *player
		players[id] = &copy
	}
	answers := make(map[string][]models.Answer, len(room.Answers))
	for id, list := range room.Answers {
		answers[id] = append([]models.Answer(nil), list...)
	}
	room.Players = players
	room.Answers = answers
	return room
}

func (memS MemoryStorage) AddPlayer(roomID string, p models.Player) error {
	memS.Mu.Lock()
	defer memS.Mu.Unlock()
	if _, ok := memS.rooms[roomID]; !ok {
		return errors.New("invalid room ID")
	}
	memS.rooms[roomID].Players[p.ID] = &p
	return nil
}

func (memS MemoryStorage) SaveAnswer(roomID string, a models.Answer) error {
	memS.Mu.Lock()
	defer memS.Mu.Unlock()
	if _, ok := memS.rooms[roomID]; !ok {
		return errors.New("invalid room ID")
	}
	memS.rooms[roomID].Answers[a.QuestionID] = append(memS.rooms[roomID].Answers[a.QuestionID], a)
	return nil
}
