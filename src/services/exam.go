package services

import (
	"fmt"
	"go-fiber-template/domain/entities"
	repo "go-fiber-template/domain/repositories"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ExamService struct {
	examRepository    repo.IExamRepository
	QuestionService   IQuestionService
	ChapterRepository repo.IChapterRepository
}

type IExamService interface {
	ExamGenerate(examRequest entities.ExamRequest, ctx *fiber.Ctx) error
	GetExamsByModuleID(moduleId string) ([]entities.ExamDataModel, error)
}

func NewExamService(examRepository repo.IExamRepository, chapterRepository repo.IChapterRepository, questionService IQuestionService) IExamService {
	return &ExamService{
		examRepository:    examRepository,
		ChapterRepository: chapterRepository,
		QuestionService:   questionService,
	}
}

func (es *ExamService) ExamGenerate(examRequest entities.ExamRequest, ctx *fiber.Ctx) error {

	exam := entities.ExamDataModel{
		ExamId:      uuid.NewString(),
		ChapterId:   examRequest.ChapterId,
		PassScore:   (examRequest.QuestionNum * 70) / 100,
		QuestionNum: examRequest.QuestionNum,
		Difficulty:  examRequest.Difficulty,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	//save exam to database
	err := es.examRepository.InsertExam(exam)
	if err != nil {
		return err
	}
	fmt.Println("Exam saved to database finished: ", exam)

	chapter, err := es.ChapterRepository.GetChaptersByChapterId(examRequest.ChapterId)
	if err != nil {
		return err
	}

	questionRequest := entities.QuestionRequest{
		Content:     chapter.ChapterContent,
		Difficulty:  examRequest.Difficulty,
		QuestionNum: examRequest.QuestionNum,
		ExamId:      exam.ExamId,
	}

	err = es.QuestionService.QuestionsCreate(questionRequest, ctx)
	if err != nil {
		return err
	}

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
