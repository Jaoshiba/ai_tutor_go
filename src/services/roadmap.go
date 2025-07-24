package services

import (
	"context"
	"fmt"
	"go-fiber-template/domain/entities"
	repo "go-fiber-template/domain/repositories"
	"mime/multipart"
	"os"
	"time"

	"github.com/google/uuid"
	"google.golang.org/genai"
)

type roadMapService struct {
	RoadMapRepo repo.IroadmapRepository
	FileService IFileService
}

type IRoadmapService interface {
	CreateRoadmap(roadmapJsonBody entities.RoadmapRequestBody, file *multipart.FileHeader) error //add userId ด้วย
}

func NewRoadmapService(roadMapRepo repo.IroadmapRepository, filrService IFileService) *roadMapService {
	return &roadMapService{
		RoadMapRepo: roadMapRepo,
		FileService: filrService,
	}
}

func (rs *roadMapService) CreateRoadmap(roadmapJsonBody entities.RoadmapRequestBody, file *multipart.FileHeader) error {

	//get content from additional file
	filetype := file.Header.Get("Content-Type")
	var content string
	if filetype == "application/pdf" {
		content, err := rs.FileService.GetPdfData(file)
		if err != nil {
			return err
		}
		content = content

	} else if filetype == "application/vnd.openxmlformats-officedocument.wordprocessingml.document" || filetype == "application/msword" {
		content, err := rs.FileService.GetDocx_DocData(file)
		if err != nil {
			fmt.Print("error docx type")
			return err
		}
		content = content

	}
	gemini_api_key := (os.Getenv("GEMINI_API_KEY"))
	if gemini_api_key == "" {
		return nil
	}
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  gemini_api_key,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return err
	}
	//send to create roadmap
	promt := fmt.Sprintf(`ฉันมีข้อมูลเบื้องต้น 3 อย่างที่ได้จากผู้ใช้:

		ชื่อ Roadmap: %s

		คำอธิบาย Roadmap: %s

		เนื้อหาที่ได้จากไฟล์แนบ: %s 
		ฉันต้องการให้คุณช่วยสร้าง Roadmap การเรียนรู้ ที่เหมาะสม โดยแบ่งเนื้อหาออกเป็น Module หรือหัวข้อหลัก ๆ ที่ควรเรียน พร้อมเรียงลำดับตามความเหมาะสมในการเรียนรู้

		ช่วยจัดรูปแบบข้อมูลให้ออกมาเป็น JSON หรือโครงสร้างที่สามารถนำไปใช้กับระบบต่อได้ (เช่นในเว็บแอป) ตัวอย่างโครงสร้างที่ต้องการ:
		{
		"title": "ชื่อของ Roadmap",
		"description": "คำอธิบาย",
		"modules": [
			{
			"title": "ชื่อ Module 1",
			"description": "คำอธิบายของ Module 1",
			"topics": [
				"หัวข้อที่ 1",
				"หัวข้อที่ 2",
				...
			]
			},
			...
		]
		}
		แต่ละ Module ควรมีหัวข้อการเรียนรู้ (topics) ที่สัมพันธ์กับเนื้อหา

		เรียงลำดับ Modules จากพื้นฐานไปขั้นสูงตามความเหมาะสม

		ไม่ต้องใส่ข้อมูลที่ไม่แน่ใจ เช่น ลิงก์ หรือไฟล์แนบ

		กรุณาสร้าง roadmap โดยอิงจากชื่อ, คำอธิบาย และเนื้อหาที่ให้ไว้ด้านบน และ **ขอเป็นภาษาไทยเป็นหลัก**`,
		roadmapJsonBody.RoadmapName, roadmapJsonBody.Description, content)

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(promt),
		nil,
	)
	if err != nil {
		fmt.Println("Error generating chapters:", err)
		return err
	}

	fmt.Println("result: ", result)

	//save roadmap to DB
	roadmap := entities.RoadmapDataModel{
		RoadmapID:   uuid.NewString(),
		RoadmapName: roadmapJsonBody.RoadmapName,
		Description: roadmapJsonBody.Description,
		Confirmed:   roadmapJsonBody.Confirmed,
		UserId:      uuid.NewString(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = rs.RoadMapRepo.InsertRoadmap(roadmap)
	if err != nil {
		fmt.Println("error insert roadmap")
		return err
	}

	return nil
}
