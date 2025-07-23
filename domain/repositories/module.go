package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"go-fiber-template/domain/entities"
)

type modulesRepository struct {
	db *sql.DB
}

type IModuleRepository interface {
	InsertModule(module entities.ModuleDataModel) error
}

func NewModulesRepository(db *sql.DB) IModuleRepository {
	return &modulesRepository{
		db: db,
	}
}

func (repo *modulesRepository) InsertModule(module entities.ModuleDataModel) error {

	fmt.Println("InsertModule called with module:", module)
	query := `
		INSERT INTO modules (
			moduleid, modulename, roadmapid, userid, createat, updateat
		) VALUES ($1, $2, $3, $4, $5, $6)`

	result, err := repo.db.ExecContext(context.Background(), query,
		module.ModuleId,
		module.ModuleName,
		module.RoadmapId,
		module.UserId,
		module.CreatedAt,
		module.UpdatedAt,
	)
	fmt.Println("result: ", result)
	if err != nil {
		return err
	}
	return nil
}
