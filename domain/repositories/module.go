package repositories

import (
	"context"
	"fmt"
	ds "go-fiber-template/domain/datasources"
	"go-fiber-template/domain/entities"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

type modulesRepository struct {
	Context    context.Context
	Collection *mongo.Collection
}

type IModuleRepository interface {
	InsertModule(module entities.ModuleDataModel) error
}

func NewModulesRepository(db *ds.MongoDB) IModuleRepository {
	return &modulesRepository{
		Context:    db.Context,
		Collection: db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("modules"),
	}
}

func (f *modulesRepository) InsertModule(module entities.ModuleDataModel) error {
	fmt.Println("InsertModule called with module:", module)
	_, err := f.Collection.InsertOne(f.Context, module)
	if err != nil {
		return err
	}
	return nil
}
