package services

import (
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
	// userID := ctx.Locals("userID").(string) // Get userId from context locals

	userID := uuid.NewString()
	if userID == "" {
		return nil, fmt.Errorf("user ID is missing from context")
	}

	course, err := rs.CourseRepo.GetCoursesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	fmt.Println("course data : ", course)

	if course == nil {
		return nil, fiber.ErrNotFound // Return a Fiber-specific error if not found
	}
	return course, nil
}

func (rs *courseService) CreateCourse(courseJsonBody entities.CourseRequestBody, file *multipart.FileHeader, fromCoures bool, ctx *fiber.Ctx) error {

	fmt.Println("Im here")

	// var serpres entities.SerpAPIResponse

	txt, err := SearchDocuments(courseJsonBody.Title, courseJsonBody.Description, ctx)
	if err != nil {
		return fmt.Errorf("failed to search documents: %w", err)
	}

	fmt.Println("Search result: ", txt)

	//จุดคั่น
	// return err

	// content := txt

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
			content = txt
		}

		ctx.Locals("content", content)

		prompt := fmt.Sprintf(`ฉันมีข้อมูลเบื้องต้น 3 อย่างที่ได้จากผู้ใช้:

			ชื่อ Course: %s

			คำอธิบาย Course: %s

			เนื้อหาจากไฟล์ที่เกียวข้อง %s

			คุณได้รับมอบหมายให้สร้างหลักสูตรการเรียนรู้ที่ครอบคลุม โดยอิงจากข้อมูลเบื้องต้นที่ผู้ใช้ให้มา ซึ่งได้แก่ ชื่อหลักสูตร: [Course Name], คำอธิบายหลักสูตร: [Course Description], และ เนื้อหาจากไฟล์ที่เกี่ยวข้อง: [Content]

			สวมบทบาทเป็นนักออกแบบหลักสูตรที่มีความรู้ความเชี่ยวชาญด้านการพัฒนาหลักสูตรและการออกแบบการเรียนการสอน เพื่อให้แน่ใจว่าเนื้อหามีการจัดระเบียบอย่างชัดเจนและมีตรรกะ

			กลุ่มเป้าหมายของคุณคือผู้สอน นักออกแบบการเรียนการสอน หรือผู้ที่ต้องการสร้างประสบการณ์การเรียนรู้ที่มีโครงสร้างสำหรับผู้เรียน

			หน้าที่ของคุณคือการสร้างโครงสร้างหลักสูตรโดยแบ่งเนื้อหาออกเป็นโมดูลหรือหัวข้อหลักที่ควรเรียนรู้ และจัดลำดับให้เหมาะสมจากระดับพื้นฐานไปจนถึงระดับสูง พร้อมทั้งระบุวัตถุประสงค์หลักของการเรียนรู้หลักสูตรนี้

			โปรดจัดรูปแบบผลลัพธ์เป็นโครงสร้าง JSON เพื่อให้ง่ายต่อการนำไปใช้งานในเว็บแอป โดยมีรูปแบบดังนี้:
			{
			"purpose": "วัตถุประสงค์ของหลักสูตรนี้คือ...",
			"modules": [
				{
				"title": "ชื่อโมดูล 1",
				"description": "คำอธิบายสำหรับโมดูล 1"
				},
				{
				"title": "ชื่อโมดูล 2",
				"description": "คำอธิบายสำหรับโมดูล 2"
				}
				]
			}
			`,
			courseJsonBody.Title, courseJsonBody.Description, content)

		modules, err := rs.GeminiService.GenerateContentFromPrompt(ctx.Context(), prompt)
		if err != nil {
			return err
		}

		courseId := uuid.NewString()
		if courseId == "" {
			fmt.Println("NUll coursei")
		}
		ctx.Locals("courseID", courseId)

		var courses entities.CourseGeminiResponse

		err = json.Unmarshal([]byte(modules), &courses)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("--- ข้อมูลหลังจาก Unmarshal ไปยัง Go struct ---")
		fmt.Printf("Course Name: %s\n", courseJsonBody.Title)
		fmt.Println("Purpose: ", courses.Purpose)
		fmt.Printf("Number of Modules: %d\n", len(courses.Modules))
		fmt.Printf("First Module Title: %s\n", courses.Modules[0].Title)
		fmt.Println("---------------------------------------------")

		fmt.Println("Module : heheheh : ", courses.Modules)

		//on web
		// userid := ctx.Locals("userID")
		userId := uuid.NewString()

		course := entities.CourseDataModel{
			CourseID:    courseId,
			Title:       courseJsonBody.Title,
			Description: courseJsonBody.Description,
			Confirmed:   courseJsonBody.Confirmed,
			UserId:      userId,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		// fmt.Println(course)

		//mocked for postman test
		// course := entities.CourseDataModel{
		// 	CourseID:    courseId,
		// 	Title:       courseJsonBody.Title,
		// 	Description: courseJsonBody.Description,
		// 	Confirmed:   courseJsonBody.Confirmed,
		// 	UserId:      userid.(string),
		// 	CreatedAt:   time.Now(),
		// 	UpdatedAt:   time.Now(),
		// }
		// ctx.Locals("userID", uuid.NewString())

		err = rs.CourseRepo.InsertCourse(course)
		if err != nil {
			fmt.Println("error insert course")
			fmt.Println(err)
			return err
		}

		// return err
		for _, moduleData := range courses.Modules {
			// fmt.Println("Module : ", moduleData)
			moduleData.Content = content
			//find title docs and insert into moduleData
			content, err := SearchDocuments(moduleData.Title, moduleData.Description, ctx)
			if err != nil {
				return fmt.Errorf("failed to search documents for module: %w", err)
			}

			moduleData.Content = content

			fmt.Println("Module Data Content: ", moduleData.Content)

			// err = rs.ModuleService.CreateModule(ctx, &moduleData)
			// if err != nil {
			// 	return err
			// }
		}
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
			return fmt.Errorf("file is nil")
		}

		courseId := uuid.NewString()
		title := file.Filename
		userid := uuid.NewString()
		course := entities.CourseDataModel{
			CourseID:    courseId,
			Title:       title,
			Description: "",
			Confirmed:   true,
			UserId:      userid,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		fmt.Println(course)

		// err := rs.CourseRepo.InsertCourse(course)
		// if err != nil {
		// 	fmt.Println("error insert course")
		// 	fmt.Println(err)
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
