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
	GetCoursesByUserID(userId string) ([]entities.CourseDataModel, error)
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

func (repo *courseRepository) GetCoursesByUserID(userId string) ([]entities.CourseDataModel, error) {
	query := `
		SELECT id, title, description, userid, createat, updateat
		FROM courses
		WHERE userid = $1
	`
	rows, err := repo.db.QueryContext(context.Background(), query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed after function returns

	var courses []entities.CourseDataModel
	for rows.Next() {
		var course entities.CourseDataModel
		if err := rows.Scan(
			&course.CourseID,
			&course.Title,
			&course.Description,
			&course.UserId,
			&course.CreatedAt,
			&course.UpdatedAt,
		); err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	if err = rows.Err(); err != nil {
		return nil, err // Check for any errors during row iteration
	}

	return courses, nil
}