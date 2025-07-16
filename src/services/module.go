package services

import (
	"fmt"
	"go-fiber-template/domain/entities"
	repo "go-fiber-template/domain/repositories"
	"mime/multipart"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ModuleService struct {
	modulesRepository repo.IModuleRepository
}

type IModuleService interface {
	CreateModule(file *multipart.FileHeader, ctx *fiber.Ctx) error
}

func NewModuleService(modulesRepository repo.IModuleRepository) IModuleService {
	return &ModuleService{
		modulesRepository: modulesRepository,
	}
}

func (ms *ModuleService) CreateModule(file *multipart.FileHeader, ctx *fiber.Ctx) error {
	filetype := file.Header.Get("Content-Type")
	fmt.Println("File header: ", filetype)

	var chapters []entities.ChapterDataModel
	err := error(nil)
	if filetype == "application/pdf" {
		chapters, err = GetPdfData(file)
		if err != nil {
			return err
		}

	} else if filetype == "application/vnd.openxmlformats-officedocument.wordprocessingml.document" || filetype == "application/msword" {
		chapters, err = GetDocx_DocData(file)
		if err != nil {
			return err
		}

	}
	// fmt.Println("Extracted chapters: ", chapters)

	exams, err := ExamGenerate(chapters)
	if err != nil {
		return err
	}

	fmt.Println("before module creation")

	module := entities.ModuleDataModel{
		ModuleId:   uuid.NewString(),
		ModuleName: file.Filename,
		RoadmapId:  "",
		Chapters:   chapters,
		Exam:       exams,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	fmt.Println("Module to be inserted: ")

	return ms.modulesRepository.InsertModule(module)
}
