package services

import (
	"context"
	"encoding/json"
	"fmt"
	"go-fiber-template/domain/entities"
	repo "go-fiber-template/domain/repositories"
	"mime/multipart"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"google.golang.org/genai"
)

type roadMapService struct {
	RoadMapRepo repo.IroadmapRepository
	FileService IFileService
}

type IRoadmapService interface {
	CreateRoadmap(roadmapJsonBody entities.RoadmapRequestBody, file *multipart.FileHeader, callFromRoadmap bool, ctx *fiber.Ctx) error //add userId ด้วย
}

func NewRoadmapService(roadMapRepo repo.IroadmapRepository) IRoadmapService {
	return &roadMapService{
		RoadMapRepo: roadMapRepo,
	}
}

func (rs *roadMapService) CreateRoadmap(roadmapJsonBody entities.RoadmapRequestBody, file *multipart.FileHeader, callFromRoadmap bool, ctx *fiber.Ctx) error {

	if callFromRoadmap {
		//get content from additional file
		roadmapName := roadmapJsonBody.RoadmapName
		roadmapDescription := roadmapJsonBody.Description
		roadmapConfirmed := roadmapJsonBody.Confirmed

		filetype := file.Header.Get("Content-Type")
		var content string
		if filetype == "application/pdf" {
			fileContent, err := GetPdfData(file, ctx)
			if err != nil {
				return err
			}
			content = fileContent

		} else if filetype == "application/vnd.openxmlformats-officedocument.wordprocessingml.document" || filetype == "application/msword" {
			fileContent, err := GetDocx_DocData(file, ctx)
			if err != nil {
				fmt.Print("error docx type")
				return err
			}
			content = fileContent

		}
		gemini_api_key := (os.Getenv("GEMINI_API_KEY"))
		if gemini_api_key == "" {
			return nil
		}

		bctx := context.Background()

		client, err := genai.NewClient(bctx, &genai.ClientConfig{
			APIKey:  gemini_api_key,
			Backend: genai.BackendGeminiAPI,
		})
		if err != nil {
			return err
		}
		fmt.Println("Creating Roadmap...")
		//send to create roadmap
		promt := fmt.Sprintf(`ฉันมีข้อมูลเบื้องต้น 3 อย่างที่ได้จากผู้ใช้:

			ชื่อ Roadmap: %s

			คำอธิบาย Roadmap: %s

			เนื้อหาที่ได้จากไฟล์แนบ: %s 
			ฉันต้องการให้คุณช่วยสร้าง Roadmap การเรียนรู้ ที่เหมาะสม โดยแบ่งเนื้อหาออกเป็น Module หรือหัวข้อหลัก ๆ ที่ควรเรียน พร้อมเรียงลำดับตามความเหมาะสมในการเรียนรู้

			ช่วยจัดรูปแบบข้อมูลให้ออกมาเป็น JSON หรือโครงสร้างที่สามารถนำไปใช้กับระบบต่อได้ (เช่นในเว็บแอป) ตัวอย่างโครงสร้างที่ต้องการ:
			{
			"modules": [
				{
				"title": "ชื่อ Module 1",
				"description": "คำอธิบายของ Module 1",
				},
				...
			]
			}
			แต่ละ Module ควรมีหัวข้อการเรียนรู้ (topics) ที่สัมพันธ์กับเนื้อหา

			เรียงลำดับ Modules จากพื้นฐานไปขั้นสูงตามความเหมาะสม

			ไม่ต้องใส่ข้อมูลที่ไม่แน่ใจ เช่น ลิงก์ หรือไฟล์แนบ

			กรุณาสร้าง roadmap โดยอิงจากชื่อ, คำอธิบาย และเนื้อหาที่ให้ไว้ด้านบน และ **ขอเป็นภาษาไทยเป็นหลัก**`,
			roadmapName, roadmapDescription, content)

		bctx = context.Background()
		result, err := client.Models.GenerateContent(
			bctx,
			"gemini-2.5-flash",
			genai.Text(promt),
			nil,
		)
		if err != nil {
			fmt.Println("Error generating chapters:", err)
			return err
		}
		roadmapString := RemoveJsonBlock(result.Text())
		if roadmapString == "" {
			return err
		}

		var geminiRes entities.RoadmapGeminiResponse
		err = json.Unmarshal([]byte(roadmapString), &geminiRes)
		if err != nil {
			return err
		}

		fmt.Println("result: ", geminiRes)

		if roadmapConfirmed {
			roadmap := entities.RoadmapDataModel{
				RoadmapID:   uuid.NewString(),
				RoadmapName: roadmapName,
				Description: roadmapDescription,
				UserId:      uuid.NewString(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			err = rs.RoadMapRepo.InsertRoadmap(roadmap)
			if err != nil {
				fmt.Println("error insert roadmap")
				return err
			}
		}

	} else {
		roadmapName := ctx.FormValue("roadmapname")
		roadmapDescription := ""

		roadmap := entities.RoadmapDataModel{
			RoadmapID:   uuid.NewString(),
			RoadmapName: roadmapName,
			Description: roadmapDescription,
			UserId:      uuid.NewString(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := rs.RoadMapRepo.InsertRoadmap(roadmap)
		if err != nil {
			fmt.Println("error insert roadmap")
			return err
		}

	}
	//save roadmap to DB

	return nil
}
