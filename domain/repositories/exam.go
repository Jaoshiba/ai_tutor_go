package repositories

import (
	"context"
	"database/sql"
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

func (db *examRepository) InsertExam(exam entities.ExamDataModel) error {

	fmt.Println("InsertExam called with exam:", exam)
	query := `
		INSERT INTO exams (
			examid, chapterid, passscore, questionnum, createdat, updatedat
		) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.db.ExecContext(context.Background(), query,
		exam.ExamId,
		exam.ChapterId,
		exam.PassScore,
		exam.QuestionNum,
		exam.CreatedAt,
		exam.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}


