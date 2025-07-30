package services

import "go-fiber-template/domain/repositories"

type PineconeService struct {
	PineconeRepository repositories.IPineconeRepository
}

type IPineconeService interface {
	EmbeedingText(text string) ([]float32, error)
}

func NewPineconeService(pineconeRepository repositories.IPineconeRepository) IPineconeService {
	return &PineconeService{
		PineconeRepository: pineconeRepository,
	}
}

func (ps *PineconeService) EmbeedingText(text string) ([]float32, error) {

	embeededValue := make([]float32, 1024)

	return embeededValue, nil
}
