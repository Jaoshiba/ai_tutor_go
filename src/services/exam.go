package services

import (
	"go-fiber-template/domain/entities"
	repo "go-fiber-template/domain/repositories"
)

type ExamService struct {
	examRepository  repo.IExamRepository
	QuestionService QuestionService
}

type IExamService interface {
	ExamGenerate(examRequest entities.ExamRequest) error
}

func NewExamService(examRepository repo.IExamRepository) IExamService {
	return &ExamService{
		examRepository: examRepository,
	}
}

func (es *ExamService) ExamGenerate(examRequest entities.ExamRequest) error {

	QuestionsCreate(examRequest.Content, examRequest.Difficulty, examRequest.QuestionNum)

	// content := examRequest.Content
	// questions, err := QuestionsCreate(content, examRequest.Difficulty, examRequest.QuestionNum)
	// if err != nil {
	// 	return err
	// }
	// exam := entities.ExamDataModel{
	// 	ExamId:      uuid.NewString(),
	// 	ModuleId:    examRequest.ModuleId,
	// 	PassScore:   (len(questions) * 70) / 100,
	// 	QuestionNum: len(questions),
	// 	CreatedAt:   time.Now(),
	// 	UpdatedAt:   time.Now(),
	// }
	// //save exam to database
	// err = es.examRepository.InsertExam(exam)
	// if err != nil {
	// 	return err
	// }

	return nil
}
