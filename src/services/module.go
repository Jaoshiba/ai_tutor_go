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
	GeminiService     IGeminiService
}

type IModuleService interface {
	CreateModule(ctx *fiber.Ctx, moduleData entities.CourseGeminiResponse, courseTitle string, courseDescription string) error
	GetModulesByCourseId(courseID string) ([]entities.ModuleDataModel, error)
	DeleteModuleByCourseId(courseID string) error
	CreateNewContent(oldcontent string, ctx *fiber.Ctx) (string, error)
}

func NewModuleService(modulesRepository repo.IModuleRepository, chapterservice IChapterService, examService IExamService, docSearchService IDocSearchService, geminiService IGeminiService) IModuleService {
	return &ModuleService{
		modulesRepository: modulesRepository,
		ChapterServices:   chapterservice,
		ExamService:       examService,
		docSearchService:  docSearchService,
		GeminiService:     geminiService,
	}
}

func (ms *ModuleService) CreateNewContent(oldcontent string, ctx *fiber.Ctx) (string, error) {

	promt := fmt.Sprintf(`คุณคืออาจารย์ผู้เชี่ยวชาญด้านการสอน  
		หน้าที่ของคุณคือการสร้างเนื้อหาใหม่จากข้อมูลอ้างอิงที่ให้ไป  
		ซึ่งเนื้อหามีที่มาหลายแหล่ง คุณต้องปฏิบัติดังนี้:

		- ห้ามคัดลอกข้อความมาโดยตรง ต้องเขียนใหม่ทั้งหมด  
		- สังเคราะห์ข้อมูลจากหลายแหล่งให้ออกมาเป็นเนื้อหาที่เป็นระบบเดียวกัน  
		- ใช้ภาษาที่เข้าใจง่าย ชัดเจน และเป็นมิตร  
		- เน้นความถูกต้อง กระชับ และเหมาะสมสำหรับใช้สอน  
		- สามารถอธิบายเสริมเพิ่มเติมเพื่อทำให้ผู้เรียนเข้าใจได้ดีขึ้น  

		เนื้อหาอ้างอิง:  
		%s

		กรุณาสร้างเนื้อหาใหม่โดยอิงจากข้อมูลด้านบน`, oldcontent)

	newcontent, err := ms.GeminiService.GenerateContentFromPrompt(ctx.Context(), promt)
	if err != nil {
		return "", err
	}
	fmt.Println("new content: ", newcontent)

	return newcontent, nil
}

func (ms *ModuleService) CreateModule(ctx *fiber.Ctx, courese entities.CourseGeminiResponse, courseTitle string, courseDescription string) error {
	// Generate a new ModuleId early, as it's needed for both the module and its chapters.

	//from coursename & description
	for i, moduleData := range courese.Modules {
		// fmt.Println("Module : ", courese)
		moduleId := uuid.NewString()
		ctx.Locals("moduleID", moduleId)
		if i > 1 {
			break
		}
		//find title docs and insert into courese
		content, err := ms.docSearchService.SearchDocuments(courseTitle, courseDescription, moduleData.Title, moduleData.Description, moduleId, ctx)
		if err != nil {
			return fmt.Errorf("failed to search documents for module: %w", err)
		}

		fmt.Printf("Module %d content: %s\n", i+1, content)

		newContent, err := ms.CreateNewContent(content, ctx)
		if err != nil {
			return err
		}

		moduleData.Content = newContent
		userIdStr := ctx.Locals("userID").(string)

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

		// --- 3. Validate courese and its Topics ---
		// if courese == nil {
		// 	fmt.Println("Error: courese is nil.")
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

		err = ms.ChapterServices.ChapterrizedText(ctx, moduleData)
		if err != nil {
			return err
		}

		// err = ms.ExamService.ExamGenerate(examRequest)
		// if err != nil {
		// 	fmt.Printf("Error generating exam for module %s: %v\n", moduleId, err)
		// 	return err // Return the error if module insertion fails.
		// }

		fmt.Println("All chapters processed for module", module) // This log should be after the loop.
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
