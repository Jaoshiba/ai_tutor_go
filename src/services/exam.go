package services

import (
	"fmt"
	"go-fiber-template/domain/entities"
	repo "go-fiber-template/domain/repositories"
	"time"

	"github.com/google/uuid"
)

type ExamService struct {
	examRepository  repo.IExamRepository
	QuestionService QuestionService
}

type IExamService interface {
	ExamGenerate(examRequest entities.ExamRequest) error
	GetExamsByModuleID(moduleId string) ([]entities.ExamDataModel, error)
}

func NewExamService(examRepository repo.IExamRepository) IExamService {
	return &ExamService{
		examRepository: examRepository,
	}
}

func (es *ExamService) ExamGenerate(examRequest entities.ExamRequest) error {

	content := examRequest.Content

	questions, err := QuestionsCreate(content, examRequest.Difficulty, examRequest.QuestionNum)
	if err != nil {
		return err
	}

	fmt.Println("questions in exam: ", questions)

	exam := entities.ExamDataModel{
		ExamId:      uuid.NewString(),
		ModuleId:    examRequest.ModuleId,
		PassScore:   (len(questions) * 70) / 100,
		QuestionNum: len(questions),
		Questions:   questions,
		RefId:       examRequest.RefId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	//save exam to database
	err = es.examRepository.InsertExam(exam)
	if err != nil {
		return err
	}
	fmt.Println("Exam saved to database finished: ", exam)

	return nil
}

func (es *ExamService) GetExamsByModuleID(moduleId string) ([]entities.ExamDataModel, error) {
	exams, err := es.examRepository.GetExamsByModuleID(moduleId)
	if err != nil {
		fmt.Println("Error getting exams from repo:", err)
		return []entities.ExamDataModel{}, err
	}

	if len(exams) == 0 {
		return []entities.ExamDataModel{}, fmt.Errorf("no exams found for module id: %s", moduleId)
	}

	return exams, nil
}

// func (es *ExamService) GetExamsByCourseID(courseId string) ([]entities.ExamDataModel, error) {

// } 