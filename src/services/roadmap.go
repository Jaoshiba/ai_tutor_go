package services

import (
	repo "go-fiber-template/domain/repositories"
	"mime/multipart"
)

type roadMapService struct {
	RoadMapRepo repo.IroadmapRepository
	FileService IFileService
}

type IRoadmapService interface {
	CreateRoadmap(roadmapName string, roadmapDescription string, file *multipart.FileHeader, confirmed bool) error //add userId ด้วย
}

func NewRoadmapService(roadMapRepo repo.IroadmapRepository) *roadMapService {
	return &roadMapService{
		RoadMapRepo: roadMapRepo,
	}
}

func (repo *roadMapService) CreateRoadmap(roadmapName string, roadmapDescription string, file *multipart.FileHeader, confirmed bool) error {

	return nil
}
