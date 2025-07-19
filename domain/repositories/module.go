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

func (db *modulesRepository) InsertModule(module entities.ModuleDataModel) error {

	fmt.Println("InsertModule called with module:", module)
	query := `
		INSERT INTO modules (
			moduleid, modulename, roadmapid, userid, createdat, updatedat
		) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.db.ExecContext(context.Background(), query,
		module.ModuleId,
		module.ModuleName,
		module.RoadmapId,
		module.UserId,
		module.CreatedAt,
		module.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
