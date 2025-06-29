package repositories

import (
	"context"
	ds "go-fiber-template/domain/datasources"
	"mime/multipart"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

type filesRepository struct {
	Context    context.Context
	Collection *mongo.Collection
}

// uploadFile implements IFileRepository.

type IFileRepository interface {
	uploadFile(file *multipart.FileHeader) (string, error)
}

func NewFilesRepository(db *ds.MongoDB) IFileRepository {
	return &filesRepository{
		Context:    db.Context,
		Collection: db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("files"),
	}
}

func (f *filesRepository) uploadFile(file *multipart.FileHeader) (string, error) {
	return "", nil
}
