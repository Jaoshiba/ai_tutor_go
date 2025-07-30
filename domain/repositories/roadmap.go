package repositories

import (
	"context"
	"database/sql"
	"go-fiber-template/domain/entities"
)

type roadMapRepository struct {
	db *sql.DB
}

type IroadmapRepository interface {
	InsertRoadmap(roadmap entities.RoadmapDataModel) error
}

func NewRoadmapRepository(db *sql.DB) IroadmapRepository {
	return &roadMapRepository{
		db: db,
	}
}

func (repo *roadMapRepository) InsertRoadmap(roadmap entities.RoadmapDataModel) error {
	query := `
	INSERT INTO courses (
		roadmapid, roadmapname, userid, createat, updateat
	) VALUES ($1, $2, $3, $4, $5)`
	_, err := repo.db.ExecContext(context.Background(), query,
		roadmap.RoadmapID,
		roadmap.RoadmapName,
		roadmap.UserId,
		roadmap.CreatedAt,
		roadmap.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
