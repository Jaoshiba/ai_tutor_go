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
	FileService       IFileService
}

type IModuleService interface {
	CreateModule(file *multipart.FileHeader, ctx *fiber.Ctx) error
}

func NewModuleService(modulesRepository repo.IModuleRepository, fileService IFileService) IModuleService {
	return &ModuleService{
		modulesRepository: modulesRepository,
		FileService:       fileService,
	}
}

func (ms *ModuleService) CreateModule(file *multipart.FileHeader, ctx *fiber.Ctx) error {
	filetype := file.Header.Get("Content-Type")
	fmt.Println("File header: ", filetype)

	fmt.Println("Extracting file content....")

	var chapters []entities.ChapterDataModel
	err := error(nil)
	if filetype == "application/pdf" {
		_, err = ms.FileService.GetPdfData(file)
		if err != nil {
			return err
		}

	} else if filetype == "application/vnd.openxmlformats-officedocument.wordprocessingml.document" || filetype == "application/msword" {
		_, err = ms.FileService.GetDocx_DocData(file)
		if err != nil {
			fmt.Print("error docx type")
			return err
		}

	}
	fmt.Println("Extracted chapters: ", chapters)

	fmt.Println("before module creation")

	module := entities.ModuleDataModel{
		ModuleId:   uuid.NewString(),
		ModuleName: file.Filename,
		RoadmapId:  uuid.NewString(),
		UserId:     uuid.NewString(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	fmt.Println("Module to be inserted: ")

	return ms.modulesRepository.InsertModule(module)
}
