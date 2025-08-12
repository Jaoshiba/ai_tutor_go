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
}

type IModuleService interface {
	CreateModule(ctx *fiber.Ctx, moduleData *entities.GenModule) error
	GetModulesByCourseID(courseID string) ([]entities.ModuleDataModel, error)
}

func NewModuleService(modulesRepository repo.IModuleRepository, chapterservice IChapterService) IModuleService {
	return &ModuleService{
		modulesRepository: modulesRepository,
		ChapterServices:   chapterservice,
	}
}

func (ms *ModuleService) CreateModule(ctx *fiber.Ctx, moduleData *entities.GenModule) error {
	// Generate a new ModuleId early, as it's needed for both the module and its chapters.
	moduleId := uuid.NewString()

	// --- 1. Safely retrieve userID from context ---
	userIdRaw := ctx.Locals("userID")
	// Always check for nil first, then perform type assertion.
	if userIdRaw == nil {
		fmt.Println("Error: User ID not found in context locals for ModuleService.")
		return fiber.NewError(fiber.StatusUnauthorized, "User ID not found in context")
	}
	userIdStr, ok := userIdRaw.(string)
	if !ok || userIdStr == "" {
		fmt.Println("Error: Invalid or missing user ID format in context locals for ModuleService.")
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or missing user ID")
	}

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
	if moduleData == nil {
		fmt.Println("Error: moduleData is nil.")
		return fmt.Errorf("module data cannot be nil")
	}

	module := entities.ModuleDataModel{
		ModuleId:    moduleId,
		ModuleName:  moduleData.Title, // Use Title from Gemini's response
		CourseId:    courseId,         // Use CourseId retrieved from context
		UserId:      userIdStr,        // Use UserId retrieved from context
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Description: moduleData.Description, // Use Description from Gemini's response
	}

	fmt.Println("Module to be inserted:", module)

	// err := ms.modulesRepository.InsertModule(module)
	// if err != nil {
	// 	fmt.Printf("Error inserting module %s into repository: %v\n", moduleId, err)
	// 	return err // Return the error if module insertion fails.
	// }
	// fmt.Println("Module successfully inserted into database.")

	err := ms.ChapterServices.ChapterrizedText(ctx, courseId, moduleData.Content)
	if err != nil {
		return err
	}

	fmt.Println("All chapters processed for module", moduleId) // This log should be after the loop.

	return nil
}
func (ms *ModuleService) GetModulesByCourseID(courseID string) ([]entities.ModuleDataModel, error) {

	modules, err := ms.modulesRepository.GetModulesByCourseID(courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve modules from repository for course %s: %w", courseID, err)
	}
	return modules, nil
}
