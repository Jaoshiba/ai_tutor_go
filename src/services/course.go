package services

import (
	// "context"
	// "encoding/json"
	"fmt"
	"go-fiber-template/domain/entities"
	repo "go-fiber-template/domain/repositories"
	"mime/multipart"
	// "os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	// "google.golang.org/genai"
)

type courseService struct {
	CourseRepo  repo.IcourseRepository
	FileService IFileService
}

type ICourseService interface {
	CreateCourse(courseJsonBody entities.CourseRequestBody, file *multipart.FileHeader, ctx *fiber.Ctx) error //add userId ด้วย
	GetCourses(ctx *fiber.Ctx) ([]entities.CourseDataModel, error)


}

func NewCourseService(courseRepo repo.IcourseRepository) ICourseService {
	return &courseService{
		CourseRepo: courseRepo,
	}
}

func (rs *courseService) GetCourses(ctx *fiber.Ctx) ([]entities.CourseDataModel, error) {
	userID := ctx.Locals("userID").(string) // Get userId from context locals

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


func (rs *courseService) CreateCourse(courseJsonBody entities.CourseRequestBody, file *multipart.FileHeader, ctx *fiber.Ctx) error {

	//get content from additional file
	fmt.Println("Im in service section nibbe")
	
	// if(file.header.Get("Content-Type")!=""){
	// 	filetype := file.Header.Get("Content-Type")
	// 	var content string
	// if filetype == "application/pdf" {
	// 	fileContent, err := GetPdfData(file, ctx) 
	// 	if err != nil {
	// 		return err
	// 	}
	// 	content = fileContent

	// } else if filetype == "application/vnd.openxmlformats-officedocument.wordprocessingml.document" || filetype == "application/msword" {
	// 	fileContent, err := GetDocx_DocData(file, ctx)
	// 	if err != nil {
	// 		fmt.Print("error docx type")
	// 		return err
	// 	}
	// 	content = fileContent
	// }
	// gemini_api_key := (os.Getenv("GEMINI_API_KEY"))
	// if gemini_api_key == "" {
	// 	return nil
	// }

	// }
	// bctx := context.Background()

	// client, err := genai.NewClient(bctx, &genai.ClientConfig{
	// 	APIKey:  gemini_api_key,
	// 	Backend: genai.BackendGeminiAPI,
	// })
	// if err != nil {
	// 	return err
	// }
	// fmt.Println("Creating Course...")
	// //send to create course
	// promt := fmt.Sprintf(`ฉันมีข้อมูลเบื้องต้น 3 อย่างที่ได้จากผู้ใช้:

	// 	ชื่อ Course: %s

	// 	คำอธิบาย Course: %s

	// 	เนื้อหาที่ได้จากไฟล์แนบ: %s 
	// 	ฉันต้องการให้คุณช่วยสร้าง Course การเรียนรู้ ที่เหมาะสม โดยแบ่งเนื้อหาออกเป็น Module หรือหัวข้อหลัก ๆ ที่ควรเรียน พร้อมเรียงลำดับตามความเหมาะสมในการเรียนรู้

	// 	ช่วยจัดรูปแบบข้อมูลให้ออกมาเป็น JSON หรือโครงสร้างที่สามารถนำไปใช้กับระบบต่อได้ (เช่นในเว็บแอป) ตัวอย่างโครงสร้างที่ต้องการ:
	// 	{
	// 	"modules": [
	// 		{
	// 		"title": "ชื่อ Module 1",
	// 		"description": "คำอธิบายของ Module 1",
	// 		},
	// 		...
	// 	]
	// 	}
	// 	แต่ละ Module ควรมีหัวข้อการเรียนรู้ (topics) ที่สัมพันธ์กับเนื้อหา

	// 	เรียงลำดับ Modules จากพื้นฐานไปขั้นสูงตามความเหมาะสม

	// 	ไม่ต้องใส่ข้อมูลที่ไม่แน่ใจ เช่น ลิงก์ หรือไฟล์แนบ

	// 	กรุณาสร้าง course โดยอิงจากชื่อ, คำอธิบาย และเนื้อหาที่ให้ไว้ด้านบน และ **ขอเป็นภาษาไทยเป็นหลัก**`,
	// 	courseJsonBody.Title, courseJsonBody.Description, content)

	// bctx = context.Background()
	// result, err := client.Models.GenerateContent(
	// 	bctx,
	// 	"gemini-2.5-flash",
	// 	genai.Text(promt),
	// 	nil,
	// )
	// if err != nil {
	// 	fmt.Println("Error generating chapters:", err)
	// 	return err
	// }
	// courseString := RemoveJsonBlock(result.Text())
	// if courseString == "" {
	// 	return err
	// }

	// var geminiRes entities.CourseGeminiResponse
	// err = json.Unmarshal([]byte(courseString), &geminiRes)
	// if err != nil {
	// 	return err
	// }

	// fmt.Println("result: ", geminiRes)

	userid := ctx.Locals("userID")

	//save course to DB
	course := entities.CourseDataModel{
		CourseID:    uuid.NewString(),
		Title:  courseJsonBody.Title,
		Description: courseJsonBody.Description,
		Confirmed:   courseJsonBody.Confirmed,
		UserId:      userid.(string),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := rs.CourseRepo.InsertCourse(course)
	if err != nil {
		fmt.Println("error insert course")
		fmt.Println(err)
		return err
	}

	return nil
}
