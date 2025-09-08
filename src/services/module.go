package services

import (
	"fmt"
	"go-fiber-template/domain/entities"
	repo "go-fiber-template/domain/repositories"

	// "mime/multipart"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ModuleService struct {
	modulesRepository repo.IModuleRepository
	ChapterServices   IChapterService
	ExamService       IExamService
	docSearchService  IDocSearchService
}

type IModuleService interface {
	CreateModule(ctx *fiber.Ctx, moduleData entities.CourseGeminiResponse, courseTitle string, courseDescription string, fromfile bool) error
	GetModulesByCourseId(courseID string) ([]entities.ModuleDataModel, error)
	DeleteModuleByCourseId(courseID string) error
}

func NewModuleService(modulesRepository repo.IModuleRepository, chapterservice IChapterService, examService IExamService, docSearchService IDocSearchService) IModuleService {
	return &ModuleService{
		modulesRepository: modulesRepository,
		ChapterServices:   chapterservice,
		ExamService:       examService,
		docSearchService:  docSearchService,
	}
}

func (ms *ModuleService) CreateModule(ctx *fiber.Ctx, moduleData entities.CourseGeminiResponse, courseTitle string, courseDescription string, fromfile bool) error {
	// Generate a new ModuleId early, as it's needed for both the module and its chapters.

	for i, moduleData := range moduleData.Modules {
		// fmt.Println("Module : ", moduleData)
		moduleId := uuid.NewString()
		ctx.Locals("moduleID", moduleId)
		if i > 5 {
			break
		}
		//find title docs and insert into moduleData
		serpReturn, err := ms.docSearchService.SearchDocuments(courseTitle, courseDescription, moduleData.Title, moduleData.Description, moduleId, ctx)
		if err != nil {
			return fmt.Errorf("failed to search documents for module: %w", err)
		}

		fmt.Printf("Module %d content: %s\n", i+1, serpReturn.Content)

		// userIdRaw := ctx.Locals("userID")
		// // Always check for nil first, then perform type assertion.
		// if userIdRaw == nil {
		// 	fmt.Println("Error: User ID not found in context locals for ModuleService.")
		// 	return fiber.NewError(fiber.StatusUnauthorized, "User ID not found in context")
		// }
		// userIdStr, ok := userIdRaw.(string)
		// if !ok || userIdStr == "" {
		// 	fmt.Println("Error: Invalid or missing user ID format in context locals for ModuleService.")
		// 	return fiber.NewError(fiber.StatusUnauthorized, "Invalid or missing user ID")
		// }
		userIdStr := uuid.NewString()

		fmt.Println("ModuleService: User ID is:", userIdStr)

		courseIdRaw := ctx.Locals("courseID")
		if courseIdRaw == nil {
			fmt.Println("Error: Course ID not found in context locals for ModuleService.")
			return fiber.NewError(fiber.StatusUnauthorized, "Course ID not found in context")
		}
		courseId, ok := courseIdRaw.(string)
		if !ok || courseId == "" {
			fmt.Println("Error: Invalid or missing course ID format in context locals for ModuleService.")
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid or missing course ID")
		}
		fmt.Println("ModuleService: Course ID is:", courseId)

		// --- 3. Validate moduleData and its Topics ---
		// if moduleData == nil {
		// 	fmt.Println("Error: moduleData is nil.")
		// 	return fmt.Errorf("module data cannot be nil")
		// }

		module := entities.ModuleDataModel{
			ModuleId:    moduleId,
			ModuleName:  moduleData.Title,
			CourseId:    courseId,
			UserId:      userIdStr,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Description: moduleData.Description,
		}

		fmt.Println("Module to be inserted:", module)

		err = ms.modulesRepository.InsertModule(module)
		if err != nil {
			fmt.Printf("Error inserting module %s into repository: %v\n", moduleId, err)
			return err // Return the error if module insertion fails.
		}
		fmt.Println("Module successfully inserted into database.")

		examRequest := entities.ExamRequest{
			ModuleId:    moduleId,
			Content:     serpReturn.Content,
			RefId:       serpReturn.RefId,
			Difficulty:  "medium",
			QuestionNum: 10,
		}
		err = ms.ExamService.ExamGenerate(examRequest)
		if err != nil {
			fmt.Printf("Error generating exam for module %s: %v\n", moduleId, err)
			return err // Return the error if module insertion fails.
		}
		// err = ms.ChapterServices.ChapterrizedText(ctx, courseId, *moduleData)
		// if err != nil {
		// 	return err
		// }

		fmt.Println("All chapters processed for module", moduleId) // This log should be after the loop.
	}

	return nil
}
func (ms *ModuleService) GetModulesByCourseId(courseID string) ([]entities.ModuleDataModel, error) {

	modules, err := ms.modulesRepository.GetModulesByCourseId(courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve modules from repository for course %s: %w", courseID, err)
	}
	return modules, nil
}

func (ms *ModuleService) DeleteModuleByCourseId(courseID string) error {

	if courseID == "" {
		return fmt.Errorf("no course id found")
	}
	err := ms.modulesRepository.DeleteModulesByCourseId(courseID)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("cant delete module")
	}

	return nil
}
