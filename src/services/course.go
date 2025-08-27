package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-fiber-template/domain/entities"
	repo "go-fiber-template/domain/repositories"
	"mime/multipart"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type courseService struct {
	CourseRepo      repo.IcourseRepository
	ModuleService   IModuleService
	GeminiService   IGeminiService
	ChapterServices IChapterService
}

type ICourseService interface {
	CreateCourse(courseJsonBody entities.CourseRequestBody, file *multipart.FileHeader, fromCoures bool, ctx *fiber.Ctx) error //add userId ด้วย
	GetCourses(ctx *fiber.Ctx) ([]entities.CourseDataModel, error)
	GetCourseDetail(ctx *fiber.Ctx, courseID string) (*entities.CourseDetailResponse, error)
	DeleteCourse(ctx *fiber.Ctx, courseID string) error
}

func NewCourseService(
	courseRepo repo.IcourseRepository,
	moduleService IModuleService, // เพิ่มเข้ามา
	geminiService IGeminiService, // เพิ่มเข้ามา
	chapterService IChapterService,
) ICourseService {
	return &courseService{
		CourseRepo:      courseRepo,
		ModuleService:   moduleService, // กำหนดค่า
		GeminiService:   geminiService, // กำหนดค่า
		ChapterServices: chapterService,
	}
}

func (rs *courseService) GetCourses(ctx *fiber.Ctx) ([]entities.CourseDataModel, error) {
	userID := ctx.Locals("userID").(string) // Get userId from context locals

	if userID == "" {
		fmt.Println("no user id")
		return nil, fmt.Errorf("user ID is missing from context")
	}

	course, err := rs.CourseRepo.GetCoursesByUserID(userID)
	if err != nil {
		fmt.Println("error get courses : ")
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	fmt.Println("course data : ", course)

	// if course == nil {
	// 	return nil, fiber.ErrNotFound // Return a Fiber-specific error if not found
	// }

	return course, nil
}

func (rs *courseService) genCourse(courseJsonBody entities.CourseRequestBody, content string, ctx context.Context) (entities.CourseGeminiResponse, error) {
	prompt := fmt.Sprintf(`ฉันมีข้อมูลเบื้องต้น 3 อย่างที่ได้จากผู้ใช้:

			ชื่อ Course: %s

			คำอธิบาย Course: %s

			เนื้อหาจากไฟล์ที่เกียวข้อง %s

			ChatGPT - Course Creation Prompt

You're tasked with creating a comprehensive learning course based on preliminary information provided by the user, including the course name, description, and relevant content.

Act as a knowledgeable course designer with expertise in curriculum development and instructional design, ensuring that the material is organized clearly and logically.

Your audience is educators, instructional designers, or anyone looking to create a structured learning experience for students.

Use the following information provided by the user: Course Name: [Course Name], Course Description: [Course Description], and Content from Related File: [Content]. Your job is to create the course structure by breaking down the content into modules or main topics that should be learned, organizing them in an appropriate sequence from basic to advanced.

Please format the output as a JSON structure for easy integration into a web app, like this example: { "modules": [ { "title": "Module Title 1", "description": "Description for Module 1", }, { "title": "Module Title 2", "description": "Description for Module 2", }, ] } Make sure your response is primarily in Thai as requested.`,
		courseJsonBody.Title, courseJsonBody.Description, content)
	modules, err := rs.GeminiService.GenerateContentFromPrompt(ctx, prompt)
	if err != nil {
		fmt.Println(err)
		return entities.CourseGeminiResponse{}, err
	}
	var courses entities.CourseGeminiResponse

	err = json.Unmarshal([]byte(modules), &courses)
	if err != nil {
		fmt.Println(err)
		return entities.CourseGeminiResponse{}, err
	}
	fmt.Println("--- ข้อมูลหลังจาก Unmarshal ไปยัง Go struct ---")
	fmt.Printf("Course Name: %s\n", courseJsonBody.Title)
	fmt.Println("Purpose: ", courses.Purpose)
	fmt.Printf("Number of Modules: %d\n", len(courses.Modules))
	fmt.Printf("First Module Title: %s\n", courses.Modules[0].Title)
	fmt.Println("---------------------------------------------")
	fmt.Println("Module : heheheh : ", courses.Modules)
	return courses, nil

}

func (rs *courseService) CreateCourse(courseJsonBody entities.CourseRequestBody, file *multipart.FileHeader, fromCoures bool, ctx *fiber.Ctx) error {

	fmt.Println("Im here")

	// var serpres entities.SerpAPIResponse

	fmt.Println("Extracting file content....")

	var content string
	if fromCoures {
		if file != nil {
			fmt.Println("Extracting file content....")
			docPath, err := SaveFileToDisk(file, ctx)
			if err != nil {
				fmt.Printf("Error saving file to disk: %v\n", err)
				return err
			}
			fileContent, err := ReadFileData(docPath, ctx)
			content = fileContent
			if err != nil {
				fmt.Printf("Error processing file with FileService: %v\n", err)
				return err
			}
		} else {
			fmt.Println("No file uploaded, skipping file processing")
			content = ""
			content = ""
		}

		ctx.Locals("content", content)

		courseId := uuid.NewString()
		if courseId == "" {
			fmt.Println("NUll coursei")
		}
		ctx.Locals("courseID", courseId)

		courses, err := rs.genCourse(courseJsonBody, content, ctx.Context())
		if err != nil {
			fmt.Println(err)
			return err
		}

		//on web
		userid := ctx.Locals("userID")

		course := entities.CourseDataModel{
			CourseID:    courseId,
			Title:       courseJsonBody.Title,
			Description: courseJsonBody.Description,
			Confirmed:   courseJsonBody.Confirmed,
			UserId:      userid.(string),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err = rs.CourseRepo.InsertCourse(course)
		if err != nil {
			fmt.Println("error insert course")
			fmt.Println(err)
			return err
		}

		// return err
		for i, moduleData := range courses.Modules {
			// fmt.Println("Module : ", moduleData)
			if i > 5 {
				break
			}
			moduleData.Content = content
			//find title docs and insert into moduleData
			content, err := SearchDocuments(courseJsonBody.Title, courseJsonBody.Description, moduleData.Title, moduleData.Description, ctx)
			if err != nil {
				return fmt.Errorf("failed to search documents for module: %w", err)
			}

			moduleData.Content = content

			fmt.Printf("Module %d content: %s\n", i+1, content)

			err = rs.ModuleService.CreateModule(ctx, &moduleData)
			if err != nil {
				return err
			}
		}

		// content, err := SearchDocuments(courseJsonBody.Title, courseJsonBody.Description, courses.Modules[0].Title, courses.Modules[0].Description, ctx)
		// if err != nil {
		// 	fmt.Println("failed to search documents for module : ", err)
		// 	return err
		// }
		// courses.Modules[0].Content = content
		// err = rs.ModuleService.CreateModule(ctx, &courses.Modules[0])
		// if err != nil {
		// 	return err
		// }

		fmt.Println("content : ", content)

	} else { //for file upload
		var content string

		if file != nil {
			fmt.Println("Extracting file content....")
			docPath, err := SaveFileToDisk(file, ctx)
			if err != nil {
				fmt.Printf("Error saving file to disk: %v\n", err)
				return err
			}

			fileContent, err := ReadFileData(docPath, ctx)
			content = fileContent
			if err != nil {
				fmt.Printf("Error processing file with FileService: %v\n", err)
				return err
			}
		} else {
			return fmt.Errorf("no file found")
		}

		courses, err := rs.genCourse(courseJsonBody, content, ctx.Context())
		if err != nil {
			return err
		}
		fmt.Println("courses : ", courses)

		courseId := uuid.NewString()
		ctx.Locals("courseID", courseId)
		ctx.Locals("content", content)
		for _, moduleData := range courses.Modules {
			// fmt.Println("Module : ", moduleData)
			moduleData.Content = content
			err = rs.ModuleService.CreateModule(ctx, &moduleData)
			if err != nil {
				return err
			}
		}

		title := file.Filename
		userid := ctx.Locals("userID")
		course := entities.CourseDataModel{
			CourseID:    courseId,
			Title:       title,
			Description: "",
			Confirmed:   true,
			UserId:      userid.(string),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		fmt.Println(course)

		err = rs.CourseRepo.InsertCourse(course)
		if err != nil {
			fmt.Println("error insert course")
			fmt.Println(err)
			return err
		}

		//create module
		// moduleData := entities.GenModule{
		// 	Title:       file.Filename,
		// 	Description: " ",
		// 	Content:     content,
		// }
		// err = rs.ModuleService.CreateModule(ctx, &moduleData)
		// if err != nil {
		// 	return err
		// }

		//create module
		moduleData := entities.GenModule{
			Title:       file.Filename,
			Description: " ",
			Content:     content,
		}
		err = rs.ModuleService.CreateModule(ctx, &moduleData)
		if err != nil {
			return err
		}

	}

	return nil
}

func (rs *courseService) GetCourseDetail(ctx *fiber.Ctx, courseID string) (*entities.CourseDetailResponse, error) {

	fmt.Println("Hello im in GetCourseDetail")
	courseData, err := rs.CourseRepo.GetCourseByID(courseID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fiber.NewError(fiber.StatusNotFound, "Course not found")
		}
		return nil, fmt.Errorf("failed to get course by ID: %w", err)
	}
	fmt.Println("after GetCourseById ")

	courseDetail := &entities.CourseDetailResponse{
		CourseID:    courseData.CourseID,
		Title:       courseData.Title,
		Description: courseData.Description,
		Confirmed:   courseData.Confirmed,
		Modules:     []entities.ModuleDetail{}, // Initialize empty slice
	}

	modulesData, err := rs.ModuleService.GetModulesByCourseID(courseID)
	fmt.Println("after GetModuleByCourseId ")
	if err != nil {
		return nil, fmt.Errorf("failed to get modules for course %s: %w", courseID, err)
	}

	for _, moduleData := range modulesData {

		chaptersData, err := rs.ChapterServices.GetChaptersByModuleID(moduleData.ModuleId)
		if err != nil {
			return nil, fmt.Errorf("failed to get chapters for module %s: %w", moduleData.ModuleId, err)
		}

		var chapterDetails []entities.ChapterDetail
		for _, chapter := range chaptersData {
			chapterDetails = append(chapterDetails, entities.ChapterDetail{
				ChapterId:      chapter.ChapterId,
				ChapterName:    chapter.ChapterName,
				ChapterContent: chapter.ChapterContent,
				IsFinished:     chapter.IsFinished,
			})
		}

		courseDetail.Modules = append(courseDetail.Modules, entities.ModuleDetail{
			ModuleId:    moduleData.ModuleId,
			ModuleName:  moduleData.ModuleName,
			Description: moduleData.Description,
			Chapters:    chapterDetails,
		})

	}

	return courseDetail, nil
}

func (rs *courseService) DeleteCourse(ctx *fiber.Ctx, courseID string) error {

	modules, err := rs.ModuleService.GetModulesByCourseID(courseID)
	if err != nil {

		return fmt.Errorf("no module with this courseid")
	}
	for _, m := range modules {
		moduleId := m.ModuleId

		fmt.Println("deleting chapters in module : ", moduleId)
		err = rs.ChapterServices.DeleteChapterByModuleID(moduleId)
		if err != nil {
			fmt.Println("Error deleting chapters for module:", moduleId, "Error:", err)
			return err
		}
	}
	err = rs.ModuleService.DeleteModuleByCourseID(courseID)
	if err != nil {
		fmt.Println("Error deleting modules for course:", courseID, "Error:", err)
		return err
	}
	err = rs.CourseRepo.DeleteCourse(courseID)
	if err!= nil {
		fmt.Println("Error deleting course:", courseID, "Error:", err)
		return err
	}
	return nil
}
