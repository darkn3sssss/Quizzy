package models

type Quiz struct {
	ID        string
	Title     string
	Questions []Question
}

type Question struct {
	ID            string
	Text          string
	Options       []string
	CorrectOption int
	TimeLimit     int
}

type Room struct {
	QuizID          string
	RoomCode        string
	CurrentQuestion int
	Status          string
	Players         []Player
	Answers         map[string]([]Answer)
}

type Player struct {
	ID    string
	Name  string
	Score int
}

type Answer struct {
	PlayerID   string
	QuestionID string
	Option     int
	AnsweredAt int64
}
