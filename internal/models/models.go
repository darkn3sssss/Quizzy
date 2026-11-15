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
	QuizID string `json:"quiz_id"`
	RoomID string `json:"room_id"`
	//CurrentQuestion int                   `json:"current_question"`
	//StartTime       int64                 `json:"start_time"`
	CountFinishPlayer int                   `json:"count_finish_player"`
	Status            string                `json:"status"`
	Players           map[string]*Player    `json:"players"`
	Answers           map[string]([]Answer) `json:"answers"`
}

type Answer struct {
	PlayerID   string `json:"players_id"`
	QuestionID string `json:"question_id"`
	Option     int    `json:"option"`
	AnsweredAt int64  `json:"answered_at"`
}

type Player struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Status          bool   `json:"status"`
	Score           int    `json:"score"`
	CurrentQuestion int    `json:"current_question"`
	StartTime       int64  `json:"start_time"`
}

type JoinRequest struct {
	Name string `json:"name"`
}

type NextQuestionRequest struct {
	QuestionIndex *int `json:"question_index"`
}

type AnswerRequest struct {
	QuestionIndex *int   `json:"question_index"`
	PlayerID      string `json:"player_id"`
	Option        int    `json:"option"`
}

type PublicQuestion struct {
	Text    string   `json:"text"`
	Options []string `json:"options"`
}

type NextQuestionResponse struct {
	// Status          string          `json:"status"`
	CurrentQuestion int             `json:"current_question"`
	Question        *PublicQuestion `json:"question,omitempty"`
}

type AnswerResponse struct {
	Correct  bool `json:"correct"`
	NewScore int  `json:"new_score"`
}

type RoomPlayerStateResponse struct {
	Question        *PublicQuestion `json:"question,omitempty"`
	CurrentQuestion int             `json:"question_index"`
	TimeRemaining   int64           `json:"time_remaining"`
	Player          Player          `json:"player"`
}

type RoomStateResponse struct {
	Status string `json:"status"`
	// Question        *PublicQuestion    `json:"question,omitempty"`
	// CurrentQuestion int                `json:"question_index"`
	// TimeRemaining   int64              `json:"time_remaining"`
	Players map[string]*Player `json:"players"`
}
