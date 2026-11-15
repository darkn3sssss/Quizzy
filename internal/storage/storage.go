package storage

import "quizzy/internal/models"

type Storage interface {
	CreateQuiz(q models.Quiz) error
	GetQuiz(id string) (models.Quiz, error)
	ListQuizzes() ([]models.Quiz, error)
	CreateRoom(r models.Room) error
	GetRoom(id string) (models.Room, error)
	UpdateRoom(id string, update func(*models.Room) error) (models.Room, error)
	AddPlayer(roomID string, p models.Player) error
	SaveAnswer(roomID string, a models.Answer) error
}
