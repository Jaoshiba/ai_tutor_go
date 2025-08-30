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
	GetExamsByModuleID(moduleId string) ([]entities.ExamDataModel, error)
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

func (er *examRepository) GetExamsByModuleID(moduleId string) ([]entities.ExamDataModel, error) {
	query := `SELECT *
	          FROM exams WHERE moduleid = $1`

	rows, err := er.db.QueryContext(context.Background(), query, moduleId)
	if err != nil {
		fmt.Println("Error querying exams:", err)
		return nil, err
	}
	defer rows.Close()

	var exams []entities.ExamDataModel

	for rows.Next() {
		var exam entities.ExamDataModel
		var questionsJSON []byte

		err := rows.Scan(
			&exam.ExamId,
			&exam.ModuleId,
			&exam.PassScore,
			&exam.QuestionNum,
			&questionsJSON,
			&exam.CreatedAt,
			&exam.UpdatedAt,
		)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}

		// แปลง JSON กลับเป็น struct slice
		err = json.Unmarshal(questionsJSON, &exam.Questions)
		if err != nil {
			fmt.Println("Error unmarshalling questions:", err)
			return nil, err
		}

		exams = append(exams, exam)
	}

	// ตรวจสอบ error หลัง loop
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return exams, nil
}

