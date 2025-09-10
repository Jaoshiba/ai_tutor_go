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
	GetExamsByModuleID(moduleId string) ([]entities.ExamDataModel, error)
}

func NewExamRepository(db *sql.DB) IExamRepository {
	return &examRepository{
		db: db,
	}
}

func (er *examRepository) InsertExam(exam entities.ExamDataModel) error {

	fmt.Println("InsertExam called with exam:", exam)

	query := `
		INSERT INTO exams (
			examid, chapterid, passscore, questionnum, difficulty, createdat, updatedat
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := er.db.ExecContext(context.Background(), query,
		exam.ExamId,
		exam.ChapterId,
		exam.PassScore,
		exam.QuestionNum,
		exam.Difficulty,
		exam.CreatedAt,
		exam.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (er *examRepository) GetExamsByRefId(refId string) ([]entities.ExamDataModel, error) {

	query := `SELECT examid, chapterid, passscore, questionnum, questions, difficulty, createdat, updatedat FROM exams WHERE refid = $1`
	rows, err := er.db.QueryContext(context.Background(), query, refId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exams []entities.ExamDataModel

	for rows.Next() {
		var exam entities.ExamDataModel
		if err := rows.Scan(
			&exam.ExamId,
			&exam.ChapterId,
			&exam.PassScore,
			&exam.QuestionNum,
			&exam.Difficulty,
			&exam.CreatedAt,
			&exam.UpdatedAt,
		); err != nil {
			return nil, err
		}
		exams = append(exams, exam)
	}

	return exams, nil

}

// GetExamsByModuleID implements IExamRepository.
func (er *examRepository) GetExamsByModuleID(moduleId string) ([]entities.ExamDataModel, error) {
	query := `SELECT examid, chapterid, passscore, questionnum, difficulty, createdat, updatedat FROM exams WHERE moduleid = $1`
	rows, err := er.db.QueryContext(context.Background(), query, moduleId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exams []entities.ExamDataModel

	for rows.Next() {
		var exam entities.ExamDataModel
		if err := rows.Scan(
			&exam.ExamId,
			&exam.ChapterId,
			&exam.PassScore,
			&exam.QuestionNum,
			&exam.Difficulty,
			&exam.CreatedAt,
			&exam.UpdatedAt,
		); err != nil {
			return nil, err
		}
		exams = append(exams, exam)
	}

	return exams, nil
}
