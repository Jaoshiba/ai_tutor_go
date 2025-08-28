package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-fiber-template/domain/entities"
)

type examRepository struct {
	db *sql.DB
}

type IExamRepository interface {
	InsertExam(exam entities.ExamDataModel) error
}

func NewExamRepository(db *sql.DB) IExamRepository {
	return &examRepository{
		db: db,
	}
}

func (er *examRepository) InsertExam(exam entities.ExamDataModel) error {

	fmt.Println("InsertExam called with exam:", exam)

	questionsJSON, err := json.Marshal(exam.Questions)
	if err != nil {
		fmt.Println("Error marshalling questions:", err)
		return err
	}

	query := `
		INSERT INTO exams (
			examid, moduleid, passscore, questionnum, questions, createdat, updatedat
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = er.db.ExecContext(context.Background(), query,
		exam.ExamId,
		exam.ModuleId,
		exam.PassScore,
		exam.QuestionNum,
		questionsJSON,
		exam.CreatedAt,
		exam.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
