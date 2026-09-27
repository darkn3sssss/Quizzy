package models

type Quiz struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Questions []Question `json:"questions"`
}

type Question struct {
	ID            string   `json:"id"`
	Text          string   `json:"text"`
	Options       []string `json:"options"`
	CorrectOption int      `json:"correct_option"`
	TimeLimit     int      `json:"time_limit"`
}

type Room struct {
	QuizID          string                `json:"quiz_id"`
	RoomID          string                `json:"room_id"`
	CurrentQuestion int                   `json:"current_question"`
	StartTime       int64                 `json:"start_time"`
	Status          string                `json:"status"`
	Players         []Player              `json:"players"`
	Answers         map[string]([]Answer) `json:"answers"`
}

type Player struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type Answer struct {
	PlayerID   string `json:"players_id"`
	QuestionID string `json:"question_id"`
	Option     int    `json:"option"`
	AnsweredAt int64  `json:"answered_at"`
}

type JoinRequest struct {
	Name string `json:"name"`
}

type PublicQuestion struct {
	Text    string   `json:"text"`
	Options []string `json:"options"`
}

type NextQuestionResponse struct {
	Status          string          `json:"status"`
	CurrentQuestion int             `json:"current_question"`
	Question        *PublicQuestion `json:"question,omitempty"`
}
