package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"go-fiber-template/domain/entities"
)

type courseRepository struct {
	db *sql.DB
}

type IcourseRepository interface {
	InsertCourse(course entities.CourseDataModel) error
	GetCoursesByUserId(userId string) ([]entities.CourseDataModel, error)
	GetCourseById(courseId string) (*entities.CourseDataModel, error)
	DeleteCourse(courseId string) error
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
		course.CourseId,
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

func (repo *courseRepository) GetCoursesByUserId(userId string) ([]entities.CourseDataModel, error) {
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
			&course.CourseId,
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
		return nil, err
	}

	return courses, nil
}

func (repo *courseRepository) GetCourseById(courseId string) (*entities.CourseDataModel, error) {
	query := `
        SELECT id, title, description, confirmed, userid, createat, updateat
        FROM courses
        WHERE id = $1
    `
	var course entities.CourseDataModel

	err := repo.db.QueryRowContext(context.Background(), query, courseId).Scan(
		&course.CourseId,
		&course.Title,
		&course.Description,
		&course.Confirmed,
		&course.UserId,
		&course.CreatedAt, // สแกนตรงเข้า time.Time
		&course.UpdatedAt, // สแกนตรงเข้า time.Time
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get course by ID %s: %w", courseId, err)
	}

	// ไม่จำเป็นต้องแปลง string เป็น time.Time แล้ว
	return &course, nil
}

// courses_repo.go
func (repo *courseRepository) DeleteCourse(courseId string) error {
	if courseId == "" {
		return fmt.Errorf("courseId is empty")
	}
	const q = `DELETE FROM courses WHERE id = $1`
	result, err := repo.db.ExecContext(context.Background(), q, courseId)
	if err != nil {
		return fmt.Errorf("delete course failed: %w", err)
	}

	if n, _ := result.RowsAffected(); n == 0 {
		return nil
	}
	fmt.Println("Deleted course:", courseId)
	return nil
}
