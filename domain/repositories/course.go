package repositories

import (
	"context"
	"database/sql"
	"go-fiber-template/domain/entities"
)

type courseRepository struct {
	db *sql.DB
}

type IcourseRepository interface {
	InsertCourse(course entities.CourseDataModel) error
}

func NewCourseRepository(db *sql.DB) IcourseRepository {
	return &courseRepository{
		db: db,
	}
}

func (repo *courseRepository) InsertCourse(course entities.CourseDataModel) error {
	query := `
	INSERT INTO courses (
		id, title, description, userid, createat, updateat
	) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := repo.db.ExecContext(context.Background(), query,
		course.CourseID,
		course.Title,
		course.Description,
		course.UserId,
		course.CreatedAt,
		course.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
