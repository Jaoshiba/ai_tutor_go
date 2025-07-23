package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"go-fiber-template/domain/entities"
)

type questionRepository struct {
	db *sql.DB
}

type IQuestionRepository interface {
	InsertQuestion(question entities.QuestionDataModel) error
}

func NewQuestionRepository(db *sql.DB) IQuestionRepository {
	return &questionRepository{
		db: db,
	}
}

func (db *questionRepository) InsertQuestion(question entities.QuestionDataModel) error {

	fmt.Println("InsertQuestion called with question:", question)
	query := `
		INSERT INTO questions (
			questionid, question, chapterid, corransid, createdat, updatedat
		) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.db.ExecContext(context.Background(), query,
		question.QuestionId,
		question.Question,
		question.ChapterId,
		question.CorrAnsId,
		question.CreateAt,
		question.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

