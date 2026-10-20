package storage

import (
	"fmt"
	"errors"
	"quizzy/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"context"
)

type DataBaseStorage struct{
	pool *pgxpool.Pool
	ctx context.Context
}

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error){
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot make pool: %w", err)
	}
	if err = pool.Ping(ctx); err != nil{
		return nil, fmt.Errorf("context: %w", err)
	}
	return pool, nil
}

func (dbs DataBaseStorage) CreateQuiz(q models.Quiz) error {
	query := "INSERT INTO quiz (title, questions) VALUES ($1, $2) RETURNING id"
	if err := dbs.pool.QueryRow(dbs.ctx, query, q.Title, q.Questions).Scan(&q.ID); err != nil{
		return fmt.Errorf("trouble with query: %w", err)
	}
	return nil
}

func (dbs DataBaseStorage) GetQuiz(id string) (models.Quiz, error) {
	var quiz models.Quiz
	query := "SELECT title, questions FROM quiz WHERE $1 = id"
	if err := dbs.pool.QueryRow(dbs.ctx, query, id).Scan(&quiz.Title, &quiz.Questions){
		return quiz, fmt.Errorf("trouble with query: %w", err)
	}
	quiz.ID = id
	return nil
}

func (dbs DataBaseStorage) ListQuizzes() ([]models.Quiz, error) {
	query := "SELECT * FROM quiz ORDER BY id"
	rows, err := dbs.pool.Query(dbs.ctx, query)
	if err != nil{
		return nil, fmt.Errorf("trouble with query: %w", err)
	}
	defer rows.Close()
	listQuizzes := []models.Quiz{}
	for rows.Next(){
		var quiz models.Quiz
		err = rows.Scan(&quiz.ID, &quiz.Title, &q.Questions)
		if err != nil{
			return nil, fmt.Errorf("cannot scan: %w", err)
		}
		listQuizzes = append(listQuizzes, quiz)
	}
	if err = rows.Err(); err != nil{
		return nil, fmt.Errorf("trouble in getting list of quizzes: %w", err)
	}
	return listQuizzes, nil
}

func (dbs DataBaseStorage) CreateRoom(r  models.Room) error{
	query := "INSERT INTO room 
	(current_question, start_time, status, players, answers)  
	VALUES ($1, $2, $3, $4, $5)
	RETURNING quiz_id, room_id"
	if err := dbs.pool.QueryRow(
		dbs.ctx, 
		query, 
		r.CurrentQuestion, 
		r.StartTime, 
		r.Status, 
		r.Players, 
		r.Answers,
	).Scan(&r.QuizID, &r.RoomID); err != nil{
		return fmt.Errorf("trouble with query: %w", err)
	}
	return nil
}

func (dbs DataBaseStorage) GetRoom(id string) (models.Room, error){
	var room models.Room
	query := "SELECT quiz_id, current_question, start_time, status, players, answers FROM room WHERE $1 = id"
	if err := dbs.pool.QueryRow(dbs.ct, query, id).Scan
	(
	&r.QuizID,
	&r.CurrentQuestion, 
	&cdr.StartTime, 
	&r.Status,
	&r.Players,
	&r.Answers,
	); err != nil{
		return room, fmt.Errorf("trouble with query: %w", err)
	}
	room.RoomID = id
	return room, nil
}
