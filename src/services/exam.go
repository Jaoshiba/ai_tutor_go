package services

import (
	"go-fiber-template/domain/entities"
	repo "go-fiber-template/domain/repositories"
	"time"

	"github.com/google/uuid"
)

type ExamService struct {
	examRepository repo.IExamRepository
	QuestionService QuestionService
}

type IExamService interface {
	ExamGenerate(chapters []entities.ChapterDataModel) error
}

func NewExamService(examRepository repo.IExamRepository) IExamService {
	return &ExamService{
		examRepository: examRepository,
	}
}

func (es *ExamService) ExamGenerate(chapters []entities.ChapterDataModel) error {

	for _, chapter := range chapters {
		content := chapter.ChapterContent
		questions, err := QuestionsCreate(content)
		if err != nil {
			return err
		}
		exam := entities.ExamDataModel{
			ExamId:      uuid.NewString(),
			ChapterId:   chapter.ChapterId,
			PassScore:   (len(questions) * 70) / 100,
			QuestionNum: len(questions),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		//save exam to database
		err = es.examRepository.InsertExam(exam)
		if err != nil {
			return err
		}
	}

	return nil
}
