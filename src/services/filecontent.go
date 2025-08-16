package services

import (
	"context"
	// "fmt"
	"time"

	"go-fiber-template/domain/entities"
	repo "go-fiber-template/domain/repositories"

	"github.com/google/uuid"
)

type IFileContentService interface {
	InsertFileContent(content, courseid, path string) error
}

type fileCotentService struct {
	repo repo.IFileContentRepository
	
}

func NewFileContentRepository( repo repo.IFileContentRepository) IFileContentService{
	return &fileCotentService{
		repo: repo, 
	}
}


func (s *fileCotentService) InsertFileContent(content, courseid, path string) error {
	newContentId := uuid.NewString()

	model := entities.FileContentModel{
		Id: newContentId,
		CourseId: courseid,
		Content: content,
		FilePath: path,
		CreatedAt: time.Now().String(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.repo.InsertFileContent(ctx, model)
}