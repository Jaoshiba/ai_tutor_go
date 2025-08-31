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
	GetModulesByCourseId(courseID string) ([]entities.ModuleDataModel, error)
	GetModuleByModuleId(moduleId string) (*entities.ModuleDataModel, error)
	DeleteModulesByCourseId(courseID string) error
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
			moduleid, modulename, courseid, userid, createat, updateat, description
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	result, err := repo.db.ExecContext(context.Background(), query,
		module.ModuleId,
		module.ModuleName,
		module.CourseId,
		module.UserId,
		module.CreatedAt,
		module.UpdatedAt,
		module.Description,
	)
	fmt.Println("result: ", result)
	if err != nil {
		return err
	}
	return nil
}

func (repo *modulesRepository) GetModuleByModuleId(moduleId string) (*entities.ModuleDataModel, error) {
	query := `
        SELECT moduleid, modulename, courseid, userid, createat, updateat, description
        FROM modules
        WHERE moduleid = $1
    `
	row := repo.db.QueryRowContext(context.Background(), query, moduleId)

	var module entities.ModuleDataModel
	if err := row.Scan(
		&module.ModuleId,
		&module.ModuleName,
		&module.CourseId,
		&module.UserId,
		&module.CreatedAt,
		&module.UpdatedAt,
		&module.Description,
	); err != nil {
		if err == sql.ErrNoRows {
			// Handle the case where no module is found
			return nil, nil
		}
		// Handle other potential errors during scanning
		return nil, fmt.Errorf("failed to scan module row for module ID %s: %w", moduleId, err)
	}

	return &module, nil
}

func (repo *modulesRepository) GetModulesByCourseId(courseID string) ([]entities.ModuleDataModel, error) {
	query := `
        SELECT moduleid, modulename, courseid, userid, createat, updateat, description
        FROM modules
        WHERE courseid = $1
        ORDER BY createat
    `
	rows, err := repo.db.QueryContext(context.Background(), query, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query modules by course ID %s: %w", courseID, err)
	}
	defer rows.Close()

	var modules []entities.ModuleDataModel
	for rows.Next() {
		var module entities.ModuleDataModel
		// สแกนตรงเข้าสู่ฟิลด์ที่เป็น time.Time ของ struct ได้เลย
		if err := rows.Scan(
			&module.ModuleId,
			&module.ModuleName,
			&module.CourseId,
			&module.UserId,
			&module.CreatedAt,
			&module.UpdatedAt,
			&module.Description,
		); err != nil {
			return nil, fmt.Errorf("failed to scan module row: %w", err)
		}

		// ไม่ต้องมีการแปลง string เป็น time.Time อีกแล้ว
		modules = append(modules, module)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during modules row iteration: %w", err)
	}

	return modules, nil
}

// modules_repo.go
func (repo *modulesRepository) DeleteModulesByCourseId(courseID string) error {
	if courseID == "" {
		return fmt.Errorf("courseID is empty")
	}
	const q = `DELETE FROM modules WHERE courseid = $1` // แก้ $1
	result, err := repo.db.ExecContext(context.Background(), q, courseID)
	if err != nil {
		return fmt.Errorf("delete modules by course_id failed: %w", err)
	}
	if n, _ := result.RowsAffected(); n > 0 {
		fmt.Printf("Deleted %d module(s) for course_id=%s\n", n, courseID)
	}
	return nil
}
