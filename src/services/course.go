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
	CourseRepo           repo.IcourseRepository
	ModuleService        IModuleService
	GeminiService        IGeminiService
	ChapterServices      IChapterService
	LearningProgressRepo repo.ILearningProgressRepository
}

type ICourseService interface {
	CreateCourse(courserequest entities.CourseRequestBody, fromfile bool, file *multipart.FileHeader, ctx *fiber.Ctx) (entities.CourseGeminiResponse, error) //add userId ด้วย
	GetCourses(ctx *fiber.Ctx) ([]entities.CourseDataModel, error)
	GetCourseDetail(ctx *fiber.Ctx, courseId string) (*entities.CourseDetailResponse, error)
	DeleteCourse(ctx *fiber.Ctx, courseId string) error
}

func NewCourseService(
	courseRepo repo.IcourseRepository,
	moduleService IModuleService, // เพิ่มเข้ามา
	geminiService IGeminiService, // เพิ่มเข้ามา
	chapterService IChapterService,
	learningprogressRepo repo.ILearningProgressRepository,
) ICourseService {
	return &courseService{
		CourseRepo:           courseRepo,
		ModuleService:        moduleService, // กำหนดค่า
		GeminiService:        geminiService, // กำหนดค่า
		ChapterServices:      chapterService,
		LearningProgressRepo: learningprogressRepo,
	}
}

func (rs *courseService) GetCourses(ctx *fiber.Ctx) ([]entities.CourseDataModel, error) {
	userID := ctx.Locals("userID").(string)

	if userID == "" {
		fmt.Println("no user id")
		return nil, fmt.Errorf("user ID is missing from context")
	}

	course, err := rs.CourseRepo.GetCoursesByUserId(userID)
	if err != nil {
		fmt.Println("error get courses : ", err)
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	fmt.Println("course data : ", course)

	// if course == nil {
	// 	return nil, fiber.ErrNotFound // Return a Fiber-specific error if not found
	// }

	return course, nil
}

func (rs *courseService) genCourse(courseJsonBody entities.CourseRequestBody, ctx context.Context) (entities.CourseGeminiResponse, error) {
	prompt := fmt.Sprintf(`ฉันมีข้อมูลเบื้องต้น 3 อย่างที่ได้จากผู้ใช้:

			ชื่อ Course: %s

			คำอธิบาย Course: %s


			ChatGPT - Course Creation Prompt

		You're tasked with creating a comprehensive learning course based on preliminary information provided by the user, including the course name, description, and relevant content.

		Act as a knowledgeable course designer with expertise in curriculum development and instructional design, ensuring that the material is organized clearly and logically.

		Your audience is educators, instructional designers, or anyone looking to create a structured learning experience for students.

		Use the following information provided by the user: Course Name: [Course Name], Course Description: [Course Description], and Content from Related File: [Content]. Your job is to create the course structure by breaking down the content into modules or main topics that should be learned, organizing them in an appropriate sequence from basic to advanced.

		Please format the output as a JSON structure for easy integration into a web app, like this example: { "modules": [ { "title": "Module Title 1", "description": "Description for Module 1", }, { "title": "Module Title 2", "description": "Description for Module 2", }, ] } Make sure your response is primarily in Thai as requested.`,
		courseJsonBody.Title, courseJsonBody.Description)

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

func (rs *courseService) RegenCourse(courseJsonBody entities.CourseRequestBody, ctx context.Context) (entities.CourseGeminiResponse, error) {
	var courses entities.CourseGeminiResponse

	promt := fmt.Sprintf(`คุณเป็นนักออกแบบหลักสูตรที่มีความเชี่ยวชาญ ได้รับมอบหมายให้สร้างโครงสร้างหลักสูตรใหม่ โดยใช้ข้อมูลที่ให้มาทั้งหมดเพื่อปรับปรุงโครงสร้างเดิมให้ดียิ่งขึ้น

		**ข้อมูลที่มี:**
		1.  **ชื่อหลักสูตร (Course Name):** "%s"
		2.  **คำอธิบายหลักสูตร (Course Description):** "%s"
		3.  **โครงสร้างหลักสูตรเดิม (Old Course Structure):**
			%s
		4.  **ความต้องการเพิ่มเติมของผู้ใช้ (User Additional Prompt):** "%s"

		**คำแนะนำสำหรับคุณ:**
		* **วิเคราะห์**โครงสร้างหลักสูตรเดิมและข้อเสนอแนะเพิ่มเติมจากผู้ใช้
		* สร้างโครงสร้างหลักสูตรใหม่ที่ **สอดคล้องกับชื่อและคำอธิบายหลักสูตร** โดยใช้ข้อมูลจาก "userAdditionalPrompt" เป็นแนวทางในการปรับปรุงจุดที่ผู้ใช้ไม่พึงพอใจในโครงสร้างเก่า
		* จัดลำดับเนื้อหาในโมดูลให้เป็นไปตามหลักการเรียนรู้จากพื้นฐานไปสู่ขั้นสูง
		* จัดรูปแบบผลลัพธ์เป็นโครงสร้าง **JSON** ดังตัวอย่าง:
			{
			"purpose": "จุดมุ่งหมายของการเรียนรู้",
			"modules": [
				{
				"title": "Module Title 1",
				"description": "Description for Module 1"
				},
				{
				"title": "Module Title 2",
				"description": "Description for Module 2"
				}
			]
			}
		* สร้างผลลัพธ์ในภาษาไทยเป็นหลัก`, courseJsonBody.Title, courseJsonBody.Description, courseJsonBody.Course, courseJsonBody.Addipromt)

	modulesFromGemini, err := rs.GeminiService.GenerateContentFromPrompt(ctx, promt)
	if err != nil {
		return courses, err
	}

	err = json.Unmarshal([]byte(modulesFromGemini), &courses)
	if err != nil {
		fmt.Println(err)
		return courses, err
	}

	return courses, nil
}

func (rs *courseService) CreateModulesFromFile(file *multipart.FileHeader, ctx *fiber.Ctx) ([]entities.GenModule, error) {
	var modules []entities.GenModule
	if file != nil {
		var content string
		fmt.Println("Extracting file content....")
		docPath, err := SaveFileToDisk(file, ctx)
		if err != nil {
			fmt.Printf("Error saving file to disk: %v\n", err)
			return modules, err
		}

		fileContent, err := ReadFileData(docPath, ctx)
		content = fileContent
		if err != nil {
			fmt.Printf("Error processing file with FileService: %v\n", err)
			return modules, err
		}

		fmt.Println("Content extracted from file:", content)

		prompt := fmt.Sprintf(`คุณคือผู้เชี่ยวชาญด้านการสร้างเนื้อหาที่สามารถจัดการเนื้อหาที่ฉันให้มาได้อย่างมีประสิทธิภาพ โดยคุณมีข้อจำกัดที่ว่า **ต้องใช้เฉพาะเนื้อหาที่ฉันให้เท่านั้น** และ **ห้ามสร้างข้อมูลหรือเนื้อหาใหม่ขึ้นมาเอง**

			หน้าที่ของคุณคือ:
			1.  **แบ่งเนื้อหา** ที่ให้มาออกเป็นส่วนๆ
			2.  สำหรับแต่ละส่วน ให้ **สร้าง object** ที่มีโครงสร้างดังต่อไปนี้:
				
					type GenModule struct {
					Title       string 
					Description string 
					Content     string 
				}
			3.  **สร้างชื่อหัวข้อ (Title)** ที่น่าสนใจและสื่อสารเนื้อหาในส่วนนั้นๆ ได้อย่างชัดเจน
			4.  **เขียนสรุปเนื้อหา (Description)** ที่กระชับและดึงดูดความสนใจผู้อ่านสำหรับส่วนนั้นๆ
			5.  **ใส่เนื้อหาต้นฉบับทั้งหมดของส่วนนั้นๆ** ลงใน Content
			6.  รวบรวม object ทั้งหมดให้อยู่ในรูป **array of objects** ในรูปแบบ JSON ที่ถูกต้อง

		**เนื้อหา:**
		%s
		`, content)

		modulesFromGemini, err := rs.GeminiService.GenerateContentFromPrompt(ctx.Context(), prompt)
		if err != nil {
			fmt.Println(err)
			return modules, err
		}

		err = json.Unmarshal([]byte(modulesFromGemini), &modules)
		if err != nil {
			fmt.Println(err)
			return modules, err
		}

		return modules, nil

	} else {
		return modules, fmt.Errorf("no file found")
	}

}

func (rs *courseService) CreateCourse(courserequest entities.CourseRequestBody, fromfile bool, file *multipart.FileHeader, ctx *fiber.Ctx) (entities.CourseGeminiResponse, error) {

	var courses entities.CourseGeminiResponse
	fmt.Println("Im here")

	fmt.Println("Extracting file content....")

	var content string
	if !fromfile {
		// if file != nil {
		// 	fmt.Println("Extracting file content....")
		// 	docPath, err := SaveFileToDisk(file, ctx)
		// 	if err != nil {
		// 		fmt.Printf("Error saving file to disk: %v\n", err)
		// 		return err
		// 	}
		// 	fileContent, err := ReadFileData(docPath, ctx)
		// 	content = fileContent
		// 	if err != nil {
		// 		fmt.Printf("Error processing file with FileService: %v\n", err)
		// 		return err
		// 	}
		// } else {
		// 	fmt.Println("No file uploaded, skipping file processing")
		// 	content = ""
		// 	content = ""
		// }
		if courserequest.Confirmed {

			fmt.Println("confirmed")

			courseId := uuid.NewString()
			ctx.Locals("courseID", courseId)
			userId := ctx.Locals("userID").(string)

			course := entities.CourseDataModel{
				CourseId:    courseId,
				Title:       courserequest.Title,
				Description: courserequest.Description,
				Confirmed:   courserequest.Confirmed,
				UserId:      userId,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			err := rs.CourseRepo.InsertCourse(course)
			if err != nil {
				fmt.Println("error insert course")
				fmt.Println(err)
				return courses, err
			}

			courses := courserequest.Course

			err = rs.ModuleService.CreateModule(ctx, courses, courserequest.Title, courserequest.Description, fromfile)
			if err != nil {
				fmt.Println("error insert module", err)
				return courses, err
			}

			fmt.Println("content : ", content)
		} else {
			if courserequest.IsFirtTime {
				fmt.Println("is first time")
				courses, err := rs.genCourse(courserequest, ctx.Context())
				if err != nil {
					fmt.Println(err)
					return courses, err
				}

				return courses, nil

			} else {
				if courserequest.Regen {
					fmt.Println("Regen courses")
					courses, err := rs.RegenCourse(courserequest, ctx.Context())
					if err != nil {
						fmt.Println(err)
						return courses, err
					}
					fmt.Println("courses from regen : ", courses)

					return courses, nil
				}
			}

		}

	} else { //for file upload

		courseId := uuid.NewString()
		ctx.Locals("courseId", courseId)

		userId := ctx.Locals("userID").(string)

		course := entities.CourseDataModel{
			CourseId:    courseId,
			Title:       courserequest.Title,
			Description: "",
			Confirmed:   true,
			UserId:      userId,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		fmt.Println(course)

		err := rs.CourseRepo.InsertCourse(course)
		if err != nil {
			fmt.Println("error insert course")
			fmt.Println(err)
			return courses, err
		}

		modules, err := rs.CreateModulesFromFile(file, ctx)
		if err != nil {
			return courses, err
		}
		fmt.Println("courses : ", modules)

		courses := entities.CourseGeminiResponse{
			Purpose: "Course created from file upload",
			Modules: modules,
		}

		err = rs.ModuleService.CreateModule(ctx, courses, courserequest.Title, courserequest.Description, fromfile)
		if err != nil {
			return courses, err
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

	}

	return courses, nil
}

func (rs *courseService) GetCourseDetail(ctx *fiber.Ctx, courseId string) (*entities.CourseDetailResponse, error) {

	fmt.Println("Hello im in GetCourseDetail")
	userIdRaw := ctx.Locals("userID")
	userId, ok := userIdRaw.(string)
	if !ok {
		fmt.Println("No user id")
	}

	courseData, err := rs.CourseRepo.GetCourseById(courseId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fiber.NewError(fiber.StatusNotFound, "Course not found")
		}
		return nil, fmt.Errorf("failed to get course by ID: %w", err)
	}
	fmt.Println("after GetCourseById ")

	var chapterDetails []entities.ChapterDetail
	var moduleDetails []entities.ModuleDetail
	var courseDetail entities.CourseDetailResponse

	modulese, err := rs.ModuleService.GetModulesByCourseId(courseId)
	if err != nil {
		return nil, fmt.Errorf("failed to get modules for course %s: %w", courseId, err)
	}

	for _, module := range modulese {

		var totalChapter int = 0
		passedChaptersId := make(map[string]bool)

		chapters, err := rs.ChapterServices.GetChaptersByModuleID(module.ModuleId)
		if err != nil {
			return nil, fmt.Errorf("failed to get chapters for module %s: %w", module.ModuleId, err)
		}
		totalChapter += len(chapters)

		passedChapter, err := rs.LearningProgressRepo.ListProgressByUser(userId)
		if err != nil {
			return nil, fmt.Errorf("failed to get learning progress for user %s: %w", userId, err)
		}

		for _, progress := range passedChapter {
			if progress.ModuleID == module.ModuleId {
				passedChaptersId[progress.ChapterID] = true
			}
		}

		for _, chapter := range chapters {
			if passedChaptersId[chapter.ChapterId] {
				chapter.Ispassed = true
			} else {
				chapter.Ispassed = false
			}
			chapterDetails = append(chapterDetails, entities.ChapterDetail{
				ChapterId:      chapter.ChapterId,
				ChapterName:    chapter.ChapterName,
				ChapterContent: chapter.ChapterContent,
				IsPassed:       chapter.Ispassed,
			})
		}
		moduleDetails = append(moduleDetails, entities.ModuleDetail{
			ModuleId:         module.ModuleId,
			ModuleName:       module.ModuleName,
			Description:      module.Description,
			Chapters:         chapterDetails,
			TotalChapters:    totalChapter,
			FinishedChapters: len(passedChaptersId),
		})
	}

	courseDetail = entities.CourseDetailResponse{
		CourseId:    courseData.CourseId,
		Title:       courseData.Title,
		Description: courseData.Description,
		Confirmed:   courseData.Confirmed,
		Modules:     moduleDetails,
	}

	return &courseDetail, nil
}

func (rs *courseService) DeleteCourse(ctx *fiber.Ctx, courseId string) error {

	modules, err := rs.ModuleService.GetModulesByCourseId(courseId)
	if err != nil {

		return fmt.Errorf("no module with this courseId")
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
	err = rs.ModuleService.DeleteModuleByCourseId(courseId)
	if err != nil {
		fmt.Println("Error deleting modules for course:", courseId, "Error:", err)
		return err
	}
	err = rs.CourseRepo.DeleteCourse(courseId)
	if err != nil {
		fmt.Println("Error deleting course:", courseId, "Error:", err)
		return err
	}
	return nil
}
