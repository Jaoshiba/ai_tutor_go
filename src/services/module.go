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
	ChapterServices   IChapterService
}

type IModuleService interface {
	CreateModule(file *multipart.FileHeader, roadmapname string, ctx *fiber.Ctx) error
}

func NewModuleService(modulesRepository repo.IModuleRepository, chapterservice IChapterService) IModuleService {
	return &ModuleService{
		modulesRepository: modulesRepository,
		ChapterServices:   chapterservice,
	}
}

func (ms *ModuleService) CreateModule(file *multipart.FileHeader, roadmapname string, ctx *fiber.Ctx) error {
	filetype := file.Header.Get("Content-Type")
	fmt.Println("File header: ", filetype)

	fmt.Println("Extracting file content....")

	var content string
	if filetype == "application/pdf" {
		fileContent, err := GetPdfData(file, ctx)
		if err != nil {
			fmt.Println("error pdf type")
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

	couresId := uuid.NewString()
	err := ms.ChapterServices.ChapterrizedText(ctx, couresId, content)
	if err != nil {
		return err
	}

	fmt.Println("before module creation")

	userIdRaw := ctx.Locals("userID")
	userIdStr, ok := userIdRaw.(string)
	if !ok || userIdStr == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or missing user ID")
	}
	fmt.Println("user id is : ", userIdStr)
	module := entities.ModuleDataModel{
		ModuleId:   couresId,
		ModuleName: file.Filename,
		RoadmapId:  uuid.NewString(),
		UserId:     userIdStr,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	fmt.Println("Module to be inserted: ")

	return ms.modulesRepository.InsertModule(module)
}
