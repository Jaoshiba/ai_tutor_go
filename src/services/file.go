package services

import (
	"fmt"
	repo "go-fiber-template/domain/repositories"
	"mime/multipart"
)

type FileService struct {
	FileRepository repo.IFileRepository
}

type IFileService interface {
	UploadFile(file *multipart.FileHeader) (string, error)
}

func NewFileService(fileRepository repo.IFileRepository) IFileService {
	return &FileService{
		FileRepository: fileRepository,
	}
}

func (f *FileService) UploadFile(file *multipart.FileHeader) (string, error) {

	fmt.Print(file)
		
	return "", nil
}
